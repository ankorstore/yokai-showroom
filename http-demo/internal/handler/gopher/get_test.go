package gopher_test

import (
	"database/sql"
	"encoding/json"
	"github.com/ankorstore/yokai-showroom/http-demo/internal/model"
	"github.com/ankorstore/yokai/fxsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankorstore/yokai-showroom/http-demo/internal"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestGetGopherHandler(t *testing.T) {
	var db *sql.DB
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(
		t,
		fxsql.RunFxSQLSeeds(),
		fx.Populate(&httpServer, &logBuffer, &traceExporter, &db),
	)

	t.Run("will 200 on success", func(t *testing.T) {
		// reset test buffers
		logBuffer.Reset()
		traceExporter.Reset()

		// [GET] /gophers/1 response assertion
		req := httptest.NewRequest(http.MethodGet, "/gophers/1", nil)
		rec := httptest.NewRecorder()
		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var gopher model.Gopher
		err := json.Unmarshal(rec.Body.Bytes(), &gopher)
		assert.NoError(t, err)

		assert.Equal(t, gopher.Name, "alice")
		assert.Equal(t, gopher.Job.String, "frontend")

		// log assertion
		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "info",
			"method":  http.MethodGet,
			"uri":     "/gophers/1",
			"status":  http.StatusOK,
			"message": "request logger",
		})

		// trace assertion
		tracetest.AssertHasTraceSpan(t,
			traceExporter,
			"GET /gophers/:id",
			semconv.HTTPMethod(http.MethodGet),
			semconv.HTTPRoute("/gophers/1"),
			semconv.HTTPStatusCode(http.StatusOK),
		)
	})

	t.Run("will 400 on invalid id", func(t *testing.T) {
		// reset test buffers
		logBuffer.Reset()
		traceExporter.Reset()

		// [GET] /gophers/invalid response assertion
		req := httptest.NewRequest(http.MethodGet, "/gophers/invalid", nil)
		rec := httptest.NewRecorder()
		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// log assertion
		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "info",
			"method":  http.MethodGet,
			"uri":     "/gophers/invalid",
			"status":  http.StatusBadRequest,
			"message": "request logger",
		})

		// trace assertion
		tracetest.AssertHasTraceSpan(t,
			traceExporter,
			"GET /gophers/:id",
			semconv.HTTPMethod(http.MethodGet),
			semconv.HTTPRoute("/gophers/invalid"),
			semconv.HTTPStatusCode(http.StatusBadRequest),
		)
	})

	t.Run("will 404 on non existing id", func(t *testing.T) {
		// reset test buffers
		logBuffer.Reset()
		traceExporter.Reset()

		// [GET] /gophers/99 response assertion
		req := httptest.NewRequest(http.MethodGet, "/gophers/99", nil)
		rec := httptest.NewRecorder()
		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		// log assertion
		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "info",
			"method":  http.MethodGet,
			"uri":     "/gophers/99",
			"status":  http.StatusNotFound,
			"message": "request logger",
		})

		// trace assertion
		tracetest.AssertHasTraceSpan(t,
			traceExporter,
			"GET /gophers/:id",
			semconv.HTTPMethod(http.MethodGet),
			semconv.HTTPRoute("/gophers/99"),
			semconv.HTTPStatusCode(http.StatusNotFound),
		)
	})

	t.Run("will 500 on internal server error", func(t *testing.T) {
		// reset test buffers
		logBuffer.Reset()
		traceExporter.Reset()

		// drop table for failure
		_, err := db.Exec("DROP TABLE gophers")
		assert.NoError(t, err)

		// [GET] /gophers/1 response assertion
		req := httptest.NewRequest(http.MethodGet, "/gophers/1", nil)
		rec := httptest.NewRecorder()
		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		// log assertion
		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "info",
			"method":  http.MethodGet,
			"uri":     "/gophers/1",
			"status":  http.StatusInternalServerError,
			"message": "request logger",
		})

		// trace assertion
		tracetest.AssertHasTraceSpan(t,
			traceExporter,
			"GET /gophers/:id",
			semconv.HTTPMethod(http.MethodGet),
			semconv.HTTPRoute("/gophers/1"),
			semconv.HTTPStatusCode(http.StatusInternalServerError),
		)
	})
}
