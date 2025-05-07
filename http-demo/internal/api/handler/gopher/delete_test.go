package gopher_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankorstore/yokai-showroom/http-demo/internal"
	"github.com/ankorstore/yokai/fxsql"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.uber.org/fx"
)

func TestDeleteGopherHandler(t *testing.T) {
	var db *sql.DB
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(
		t,
		fxsql.RunFxSQLSeeds(),
		fx.Populate(&httpServer, &logBuffer, &traceExporter, &db),
	)

	t.Run("should 204 on success", func(t *testing.T) {
		// reset test buffers
		logBuffer.Reset()
		traceExporter.Reset()

		// [DELETE] /gophers/1 response assertion
		req := httptest.NewRequest(http.MethodDelete, "/gophers/1", nil)
		rec := httptest.NewRecorder()
		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)

		// log assertion
		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "info",
			"method":  http.MethodDelete,
			"uri":     "/gophers/1",
			"status":  http.StatusNoContent,
			"message": "request logger",
		})

		// trace assertion
		tracetest.AssertHasTraceSpan(t,
			traceExporter,
			"DELETE /gophers/:id",
			semconv.HTTPMethod(http.MethodDelete),
			semconv.HTTPRoute("/gophers/1"),
			semconv.HTTPStatusCode(http.StatusNoContent),
		)
	})

	t.Run("should 400 on invalid id", func(t *testing.T) {
		// reset test buffers
		logBuffer.Reset()
		traceExporter.Reset()

		// [DELETE] /gophers/invalid response assertion
		req := httptest.NewRequest(http.MethodDelete, "/gophers/invalid", nil)
		rec := httptest.NewRecorder()
		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// log assertion
		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "info",
			"method":  http.MethodDelete,
			"uri":     "/gophers/invalid",
			"status":  http.StatusBadRequest,
			"message": "request logger",
		})

		// trace assertion
		tracetest.AssertHasTraceSpan(t,
			traceExporter,
			"DELETE /gophers/:id",
			semconv.HTTPMethod(http.MethodDelete),
			semconv.HTTPRoute("/gophers/invalid"),
			semconv.HTTPStatusCode(http.StatusBadRequest),
		)
	})

	t.Run("should 404 on non existing id", func(t *testing.T) {
		// reset test buffers
		logBuffer.Reset()
		traceExporter.Reset()

		// [DELETE] /gophers/99 response assertion
		req := httptest.NewRequest(http.MethodDelete, "/gophers/99", nil)
		rec := httptest.NewRecorder()
		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		// log assertion
		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "info",
			"method":  http.MethodDelete,
			"uri":     "/gophers/99",
			"status":  http.StatusNotFound,
			"message": "request logger",
		})

		// trace assertion
		tracetest.AssertHasTraceSpan(t,
			traceExporter,
			"DELETE /gophers/:id",
			semconv.HTTPMethod(http.MethodDelete),
			semconv.HTTPRoute("/gophers/99"),
			semconv.HTTPStatusCode(http.StatusNotFound),
		)
	})

	t.Run("should 500 on internal server error", func(t *testing.T) {
		// reset test buffers
		logBuffer.Reset()
		traceExporter.Reset()

		// drop table for failure
		_, err := db.Exec("DROP TABLE gophers")
		assert.NoError(t, err)

		// [DELETE] /gophers/1 response assertion
		req := httptest.NewRequest(http.MethodDelete, "/gophers/1", nil)
		rec := httptest.NewRecorder()
		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		// log assertion
		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "info",
			"method":  http.MethodDelete,
			"uri":     "/gophers/1",
			"status":  http.StatusInternalServerError,
			"message": "request logger",
		})

		// trace assertion
		tracetest.AssertHasTraceSpan(t,
			traceExporter,
			"DELETE /gophers/:id",
			semconv.HTTPMethod(http.MethodDelete),
			semconv.HTTPRoute("/gophers/1"),
			semconv.HTTPStatusCode(http.StatusInternalServerError),
		)
	})
}
