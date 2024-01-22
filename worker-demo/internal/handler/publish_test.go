package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ankorstore/yokai-showroom/worker-demo/internal"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestPublishHandlerSuccess(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", fmt.Sprintf("%s/configs", internal.RootDir))
	t.Setenv("PUBSUB_PROJECT_ID", "worker-demo-project")

	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter
	var metricsRegistry *prometheus.Registry

	app := internal.Bootstrapper.BootstrapTestApp(
		t,
		fx.Populate(
			&httpServer,
			&logBuffer,
			&traceExporter,
			&metricsRegistry,
		),
	)

	app.RequireStart()

	time.Sleep(1 * time.Second)

	req := httptest.NewRequest(http.MethodGet, "/publish?message=test", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	app.RequireStop()

	// response assertion
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "publication of message test success.")

	ll, _ := logBuffer.Records()
	for _, l := range ll {
		fmt.Printf("%v\n", l)
	}

	// logs assertion
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"service": "worker-demo-app",
		"method":  "GET",
		"uri":     "/publish?message=test",
	})

	// trace assertion
	tracetest.AssertHasTraceSpan(t, traceExporter, "GET /publish")

	// metrics assertion
	expectedMetric := `
		# HELP worker_demo_app_messages_published_total Total number of published messages
		# TYPE worker_demo_app_messages_published_total counter
		worker_demo_app_messages_published_total 1
	`

	err := testutil.GatherAndCompare(
		metricsRegistry,
		strings.NewReader(expectedMetric),
		"worker_demo_app_messages_published_total",
	)
	assert.NoError(t, err)
}
