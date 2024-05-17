package gopher_test

import (
	"database/sql"
	"encoding/json"
	"github.com/ankorstore/yokai-showroom/http-demo/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankorstore/yokai-showroom/http-demo/internal"
	"github.com/ankorstore/yokai-showroom/http-demo/internal/repository"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestGetGopherHandlerSuccess(t *testing.T) {
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter
	var repository *repository.GopherRepository

	internal.RunTest(t, fx.Populate(&httpServer, &logBuffer, &traceExporter, &repository))

	// [GET] /gophers/1 response assertion
	req := httptest.NewRequest(http.MethodGet, "/gophers/1", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var gopher model.Gopher
	err := json.Unmarshal(rec.Body.Bytes(), &gopher)
	assert.NoError(t, err)

	assert.Equal(t, gopher.Name, "alice")
	assert.Equal(t, gopher.Job, "architect")

	// log assertion
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "called get gopher",
	})

	// trace assertion
	tracetest.AssertHasTraceSpan(t, traceExporter, "get gopher service")
}

func TestGetGopherHandlerBadRequestErrorOnInvalidId(t *testing.T) {
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(t, fx.Populate(&httpServer, &logBuffer, &traceExporter))

	// [GET] /gophers/invalid response assertion
	req := httptest.NewRequest(http.MethodGet, "/gophers/invalid", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid gopher id")
}

func TestGetGopherHandlerNotFoundErrorOnMissingId(t *testing.T) {
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(t, fx.Populate(&httpServer, &logBuffer, &traceExporter))

	// [GET] /gophers/1 response assertion, database not populated
	req := httptest.NewRequest(http.MethodGet, "/gophers/1", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "cannot get gopher with id 1")
}

func TestGetGopherHandlerInternalServerErrorOnMissingTable(t *testing.T) {
	var db *sql.DB
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(t, fx.Populate(&db, &httpServer, &logBuffer, &traceExporter))

	// drop table for failure
	_, err := db.Exec("DROP TABLE gophers")
	assert.NoError(t, err)

	// [GET] /gophers/1 response assertion
	req := httptest.NewRequest(http.MethodGet, "/gophers/1", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "cannot get gopher")
}
