package worker_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ankorstore/yokai-contrib/fxgcppubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/reactor/ack"
	"github.com/ankorstore/yokai-showroom/worker-demo/internal"
	"github.com/ankorstore/yokai/fxconfig"
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

	var publisher fxgcppubsub.Publisher
	var supervisor ack.AckSupervisor
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter
	var metricsRegistry *prometheus.Registry

	// bootstrap test app
	app := internal.Bootstrapper.BootstrapTestApp(
		t,
		// config lookup
		fxconfig.AsConfigPath(fmt.Sprintf("%s/configs/", internal.RootDir)),
		// pub/sub topic and subscription preparation
		fxgcppubsub.PrepareTopicAndSubscription(fxgcppubsub.PrepareTopicAndSubscriptionParams{
			TopicID:        "test-topic",
			SubscriptionID: "test-subscription",
		}),
		// populate
		fx.Populate(&publisher, &supervisor, &logBuffer, &traceExporter, &metricsRegistry),
	)

	// start test app
	app.RequireStart()

	// subscription message ack waiter
	waiter := supervisor.StartNackWaiter("test-subscription")

	// publish test message
	testMessage := "test message"

	result, err := publisher.Publish(ctx, "test-topic", testMessage)
	assert.NoError(t, err)

	id, err := result.Get(ctx)
	assert.NoError(t, err)

	// wait for ack for max 1 second
	_, err = waiter.WaitMaxDuration(ctx, time.Second)
	assert.NoError(t, err)

	// stop test app
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
		"subscribe-worker message",
		attribute.String("Worker", "subscribe-worker"),
		attribute.String("Message", testMessage),
	)

	// metrics assertion
	expectedMetric := `
		# HELP subscriber_messages_received_total Total number of received messages
		# TYPE subscriber_messages_received_total counter
		subscriber_messages_received_total 1
	`

	err = testutil.GatherAndCompare(
		metricsRegistry,
		strings.NewReader(expectedMetric),
		"subscriber_messages_received_total",
	)
	assert.NoError(t, err)
}
