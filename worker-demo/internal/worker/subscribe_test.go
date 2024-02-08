package worker_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai-showroom/worker-demo/internal"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/fx"
)

func TestSubscribeWorker(t *testing.T) {
	ctx := context.Background()

	testMessage := "test message"

	// env vars
	t.Setenv("APP_CONFIG_PATH", fmt.Sprintf("%s/configs", internal.RootDir))

	var client *pubsub.Client
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter
	var metricsRegistry *prometheus.Registry

	// bootstrap test app
	app := internal.Bootstrapper.BootstrapTestApp(
		t,
		fx.Populate(
			&client,
			&logBuffer,
			&traceExporter,
			&metricsRegistry,
		),
	)

	// start test app
	app.RequireStart()

	// publish test message
	result := client.Topic("test-topic").Publish(ctx, &pubsub.Message{
		Data: []byte(testMessage),
	})

	id, err := result.Get(ctx)
	assert.NoError(t, err)

	// stop test app (after 1 sec wait to avoid test flakiness)
	time.Sleep(1 * time.Second)

	app.RequireStop()

	// logs assertion
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"service": "worker-demo-app",
		"module":  "worker",
		"worker":  "subscribe-worker",
		"message": fmt.Sprintf("received message: id=%v, data=%s", id, testMessage),
	})

	// trace assertion
	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"subscribe-worker span",
		attribute.String("Worker", "subscribe-worker"),
	)

	// metrics assertion
	expectedMetric := `
		# HELP messages_received_total Total number of received messages
		# TYPE messages_received_total counter
		messages_received_total 1
	`

	err = testutil.GatherAndCompare(
		metricsRegistry,
		strings.NewReader(expectedMetric),
		"messages_received_total",
	)
	assert.NoError(t, err)
}
