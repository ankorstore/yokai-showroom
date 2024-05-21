package gopher_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankorstore/yokai-showroom/http-demo/internal"
	"github.com/ankorstore/yokai-showroom/http-demo/internal/model"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestCreateGopherHandler(t *testing.T) {
	var db *sql.DB
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(t, fx.Populate(&httpServer, &logBuffer, &traceExporter, &db))

	t.Run("will 201 on success", func(t *testing.T) {
		// reset test buffers
		logBuffer.Reset()
		traceExporter.Reset()

		// [POST] /gophers response assertion
		data := `{"name": "test", "job": "test"}`
		req := httptest.NewRequest(http.MethodPost, "/gophers", bytes.NewBuffer([]byte(data)))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var gopher model.Gopher
		err := json.Unmarshal(rec.Body.Bytes(), &gopher)
		assert.NoError(t, err)

		assert.Equal(t, gopher.Name, "test")
		assert.Equal(t, gopher.Job.String, "test")

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
	})

	t.Run("will 400 on invalid data", func(t *testing.T) {
		// reset test buffers
		logBuffer.Reset()
		traceExporter.Reset()

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
	})

	t.Run("will 500 on internal server error", func(t *testing.T) {
		// reset test buffers
		logBuffer.Reset()
		traceExporter.Reset()

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
	})
}
