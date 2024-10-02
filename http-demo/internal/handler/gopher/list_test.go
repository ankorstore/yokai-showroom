package gopher_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankorstore/yokai-showroom/http-demo/internal"
	"github.com/ankorstore/yokai-showroom/http-demo/internal/model"
	"github.com/ankorstore/yokai/fxsql"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.uber.org/fx"
)

func TestListGopherHandler(t *testing.T) {
	var db *sql.DB
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(
		t,
		fxsql.RunFxSQLSeeds(),
		fx.Populate(&httpServer, &logBuffer, &traceExporter, &db),
	)

	t.Run("should 200 on success without filter", func(t *testing.T) {
		// reset test buffers
		logBuffer.Reset()
		traceExporter.Reset()

		// [GET] /gophers response assertion
		req := httptest.NewRequest(http.MethodGet, "/gophers", nil)
		rec := httptest.NewRecorder()
		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var gophers []model.Gopher
		err := json.Unmarshal(rec.Body.Bytes(), &gophers)
		assert.NoError(t, err)

		assert.Len(t, gophers, 5)
		assert.Equal(t, gophers[0].Name, "alice")
		assert.Equal(t, gophers[0].Job.String, "frontend")
		assert.Equal(t, gophers[1].Name, "bob")
		assert.Equal(t, gophers[1].Job.String, "backend")
		assert.Equal(t, gophers[2].Name, "carl")
		assert.Equal(t, gophers[2].Job.String, "backend")
		assert.Equal(t, gophers[3].Name, "dan")
		assert.Equal(t, gophers[3].Job.String, "frontend")
		assert.Equal(t, gophers[4].Name, "elvis")
		assert.Equal(t, gophers[4].Job.String, "backend")

		// log assertion
		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "info",
			"method":  http.MethodGet,
			"uri":     "/gophers",
			"status":  http.StatusOK,
			"message": "request logger",
		})

		// trace assertion
		tracetest.AssertHasTraceSpan(t,
			traceExporter,
			"GET /gophers",
			semconv.HTTPMethod(http.MethodGet),
			semconv.HTTPRoute("/gophers"),
			semconv.HTTPStatusCode(http.StatusOK),
		)
	})

	t.Run("should 200 on success with name filter", func(t *testing.T) {
		// reset test buffers
		logBuffer.Reset()
		traceExporter.Reset()

		// [GET] /gophers response assertion
		req := httptest.NewRequest(http.MethodGet, "/gophers?name=carl", nil)
		rec := httptest.NewRecorder()
		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var gophers []model.Gopher
		err := json.Unmarshal(rec.Body.Bytes(), &gophers)
		assert.NoError(t, err)

		assert.Len(t, gophers, 1)
		assert.Equal(t, gophers[0].Name, "carl")
		assert.Equal(t, gophers[0].Job.String, "backend")

		// log assertion
		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "info",
			"method":  http.MethodGet,
			"uri":     "/gophers?name=carl",
			"status":  http.StatusOK,
			"message": "request logger",
		})

		// trace assertion
		tracetest.AssertHasTraceSpan(t,
			traceExporter,
			"GET /gophers",
			semconv.HTTPMethod(http.MethodGet),
			semconv.HTTPRoute("/gophers"),
			semconv.HTTPStatusCode(http.StatusOK),
		)
	})

	t.Run("should 200 on success with job filter", func(t *testing.T) {
		// reset test buffers
		logBuffer.Reset()
		traceExporter.Reset()

		// [GET] /gophers response assertion
		req := httptest.NewRequest(http.MethodGet, "/gophers?job=backend", nil)
		rec := httptest.NewRecorder()
		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var gophers []model.Gopher
		err := json.Unmarshal(rec.Body.Bytes(), &gophers)
		assert.NoError(t, err)

		assert.Len(t, gophers, 3)
		assert.Equal(t, gophers[0].Name, "bob")
		assert.Equal(t, gophers[0].Job.String, "backend")
		assert.Equal(t, gophers[1].Name, "carl")
		assert.Equal(t, gophers[1].Job.String, "backend")
		assert.Equal(t, gophers[2].Name, "elvis")
		assert.Equal(t, gophers[2].Job.String, "backend")

		// log assertion
		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "info",
			"method":  http.MethodGet,
			"uri":     "/gophers?job=backend",
			"status":  http.StatusOK,
			"message": "request logger",
		})

		// trace assertion
		tracetest.AssertHasTraceSpan(t,
			traceExporter,
			"GET /gophers",
			semconv.HTTPMethod(http.MethodGet),
			semconv.HTTPRoute("/gophers"),
			semconv.HTTPStatusCode(http.StatusOK),
		)
	})

	t.Run("should 400 on invalid id", func(t *testing.T) {
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

	t.Run("should 500 on internal server error", func(t *testing.T) {
		// reset test buffers
		logBuffer.Reset()
		traceExporter.Reset()

		// drop table for failure
		_, err := db.Exec("DROP TABLE gophers")
		assert.NoError(t, err)

		// [GET] /gophers/1 response assertion
		req := httptest.NewRequest(http.MethodGet, "/gophers", nil)
		rec := httptest.NewRecorder()
		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		// log assertion
		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":   "info",
			"method":  http.MethodGet,
			"uri":     "/gophers",
			"status":  http.StatusInternalServerError,
			"message": "request logger",
		})

		// trace assertion
		tracetest.AssertHasTraceSpan(t,
			traceExporter,
			"GET /gophers",
			semconv.HTTPMethod(http.MethodGet),
			semconv.HTTPRoute("/gophers"),
			semconv.HTTPStatusCode(http.StatusInternalServerError),
		)
	})
}
