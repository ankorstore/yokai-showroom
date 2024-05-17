package gopher_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankorstore/yokai-showroom/http-demo/db/sqlc"
	"github.com/ankorstore/yokai-showroom/http-demo/internal"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestCreateGopherHandlerSuccess(t *testing.T) {
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(t, fx.Populate(&httpServer, &logBuffer, &traceExporter))

	// [POST] /gophers response assertion
	data := `{"name": "test", "job": "test"}`
	req := httptest.NewRequest(http.MethodPost, "/gophers", bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var gopher sqlc.Gopher
	err := json.Unmarshal(rec.Body.Bytes(), &gopher)
	assert.NoError(t, err)

	assert.Equal(t, gopher.Name, "test")
	assert.Equal(t, gopher.Job, "test")

	// log assertion
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "called create gopher",
	})

	// trace assertion
	tracetest.AssertHasTraceSpan(t, traceExporter, "create gopher service")
}

func TestCreateGopherHandlerBadRequestErrorOnInvalidData(t *testing.T) {
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(t, fx.Populate(&httpServer, &logBuffer, &traceExporter))

	// [POST] /gophers response assertion
	data := `invalid`
	req := httptest.NewRequest(http.MethodPost, "/gophers", bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid parameters")
}

func TestCreateGopherHandlerInternalServerErrorOnMissingTable(t *testing.T) {
	var db *sql.DB
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(t, fx.Populate(&db, &httpServer, &logBuffer, &traceExporter))

	// drop table for failure
	_, err := db.Exec("DROP TABLE gophers")
	assert.NoError(t, err)

	// [POST] /gophers response assertion
	req := httptest.NewRequest(http.MethodPost, "/gophers", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "cannot create gopher")
}
