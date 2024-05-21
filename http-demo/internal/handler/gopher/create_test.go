package gopher_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/ankorstore/yokai-showroom/http-demo/internal"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.uber.org/fx"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ankorstore/yokai-showroom/http-demo/internal/model"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateGopherHandler(t *testing.T) {
	var db *sql.DB
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter
	var metricsRegistry *prometheus.Registry

	t.Run("should 201 on success", func(t *testing.T) {
		// run test
		internal.RunTest(t, fx.Populate(&db, &httpServer, &logBuffer, &traceExporter, &metricsRegistry))

		// [POST] /gophers response assertion
		data := `{"name": "test name", "job": "test job"}`
		req := httptest.NewRequest(http.MethodPost, "/gophers", bytes.NewBuffer([]byte(data)))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var gopher model.Gopher
		err := json.Unmarshal(rec.Body.Bytes(), &gopher)
		assert.NoError(t, err)

		assert.Equal(t, gopher.Name, "test name")
		assert.Equal(t, gopher.Job.String, "test job")

		// log assertion
		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "info",
			"method":  http.MethodPost,
			"uri":     "/gophers",
			"status":  http.StatusCreated,
			"message": "request logger",
		})

		// trace assertion
		tracetest.AssertHasTraceSpan(t,
			traceExporter,
			"POST /gophers",
			semconv.HTTPMethod(http.MethodPost),
			semconv.HTTPRoute("/gophers"),
			semconv.HTTPStatusCode(http.StatusCreated),
		)

		// metrics assertion
		expectedMetric := `
		# HELP http_server_requests_total Number of processed HTTP requests
		# TYPE http_server_requests_total counter
		http_server_requests_total{method="POST",path="/gophers",status="2xx"} 1
	`
		err = testutil.GatherAndCompare(
			metricsRegistry,
			strings.NewReader(expectedMetric),
			"http_server_requests_total",
		)
		assert.NoError(t, err)
	})

	t.Run("should 400 on invalid data", func(t *testing.T) {
		// run test
		internal.RunTest(t, fx.Populate(&db, &httpServer, &logBuffer, &traceExporter, &metricsRegistry))

		// [POST] /gophers response assertion
		data := `invalid`
		req := httptest.NewRequest(http.MethodPost, "/gophers", bytes.NewBuffer([]byte(data)))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid parameters")

		// log assertion
		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "info",
			"method":  http.MethodPost,
			"uri":     "/gophers",
			"status":  http.StatusBadRequest,
			"message": "request logger",
		})

		// trace assertion
		tracetest.AssertHasTraceSpan(t,
			traceExporter,
			"POST /gophers",
			semconv.HTTPMethod(http.MethodPost),
			semconv.HTTPRoute("/gophers"),
			semconv.HTTPStatusCode(http.StatusBadRequest),
		)

		// metrics assertion
		expectedMetric := `
		# HELP http_server_requests_total Number of processed HTTP requests
		# TYPE http_server_requests_total counter
		http_server_requests_total{method="POST",path="/gophers",status="4xx"} 1
	`
		err := testutil.GatherAndCompare(
			metricsRegistry,
			strings.NewReader(expectedMetric),
			"http_server_requests_total",
		)
		assert.NoError(t, err)
	})

	t.Run("should 500 on internal server error", func(t *testing.T) {
		// run test
		internal.RunTest(t, fx.Populate(&db, &httpServer, &logBuffer, &traceExporter, &metricsRegistry))

		// drop table for failure
		_, err := db.Exec("DROP TABLE gophers")
		assert.NoError(t, err)

		// [POST] /gophers response assertion
		req := httptest.NewRequest(http.MethodPost, "/gophers", nil)
		rec := httptest.NewRecorder()
		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "cannot create gopher")

		// log assertion
		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "info",
			"method":  http.MethodPost,
			"uri":     "/gophers",
			"status":  http.StatusInternalServerError,
			"message": "request logger",
		})

		// trace assertion
		tracetest.AssertHasTraceSpan(t,
			traceExporter,
			"POST /gophers",
			semconv.HTTPMethod(http.MethodPost),
			semconv.HTTPRoute("/gophers"),
			semconv.HTTPStatusCode(http.StatusInternalServerError),
		)

		// metrics assertion
		expectedMetric := `
		# HELP http_server_requests_total Number of processed HTTP requests
		# TYPE http_server_requests_total counter
		http_server_requests_total{method="POST",path="/gophers",status="5xx"} 1
	`
		err = testutil.GatherAndCompare(
			metricsRegistry,
			strings.NewReader(expectedMetric),
			"http_server_requests_total",
		)
		assert.NoError(t, err)
	})
}
