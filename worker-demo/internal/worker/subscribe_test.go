package worker_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ankorstore/yokai-showroom/worker-demo/internal"
	"github.com/ankorstore/yokai-showroom/worker-demo/internal/service"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestSubscribeWorker(t *testing.T) {
	// env vars
	t.Setenv("APP_CONFIG_PATH", fmt.Sprintf("%s/configs", internal.RootDir))
	t.Setenv("PUBSUB_PROJECT_ID", "worker-demo-project")

	var publisher service.Publisher
	var logBuffer logtest.TestLogBuffer
	var metricsRegistry *prometheus.Registry

	// test app
	app := internal.Bootstrapper.BootstrapTestApp(
		t,
		fx.Populate(
			&publisher,
			&logBuffer,
			&metricsRegistry,
		),
	)

	// start test app
	app.RequireStart()

	// publish test message, with some sleep around to avoid flaky test (local gRPC conn)
	time.Sleep(1 * time.Second)

	err := publisher.Publish(context.Background(), "test")
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)

	// stop test app
	app.RequireStop()

	// logs assertion
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"service": "worker-demo-app",
		"module":  "worker",
		"worker":  "subscribe-worker",
		"message": "received new message: test",
	})

	// metrics assertion
	expectedMetric := `
		# HELP worker_demo_app_messages_received_total Total number of received messages
        # TYPE worker_demo_app_messages_received_total counter
        worker_demo_app_messages_received_total 1
	`

	err = testutil.GatherAndCompare(
		metricsRegistry,
		strings.NewReader(expectedMetric),
		"worker_demo_app_messages_received_total",
	)
	assert.NoError(t, err)
}
