package gopher_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankorstore/yokai-showroom/http-demo/db/seeds"
	"github.com/ankorstore/yokai-showroom/http-demo/internal"
	"github.com/ankorstore/yokai/fxsql"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestDeleteGopherHandlerSuccess(t *testing.T) {
	var httpServer *echo.Echo

	internal.RunTest(
		t,
		fxsql.AsSQLSeeds(seeds.NewGophersSeed),
		fxsql.RunFxSQLSeeds(),
		fx.Populate(&httpServer),
	)

	// [DELETE] /gophers/1 response assertion
	req := httptest.NewRequest(http.MethodDelete, "/gophers/1", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)

	// log assertion

}

func TestDeleteGopherHandlerBadRequestErrorOnInvalidId(t *testing.T) {
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(t, fx.Populate(&httpServer, &logBuffer, &traceExporter))

	// [DELETE] /gophers/invalid response assertion
	req := httptest.NewRequest(http.MethodDelete, "/gophers/invalid", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid gopher id")
}

func TestDeleteGopherHandlerNotFoundErrorOnMissingId(t *testing.T) {
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(t, fx.Populate(&httpServer, &logBuffer, &traceExporter))

	// [DELETE] /gophers/1 response assertion, database not populated
	req := httptest.NewRequest(http.MethodDelete, "/gophers/1", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "cannot get gopher with id 1")
}

func TestDeleteGopherHandlerInternalServerErrorOnMissingTable(t *testing.T) {
	var db *sql.DB
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(t, fx.Populate(&db, &httpServer, &logBuffer, &traceExporter))

	// drop table for failure
	_, err := db.Exec("DROP TABLE gophers")
	assert.NoError(t, err)

	// [DELETE] /gophers/1 response assertion
	req := httptest.NewRequest(http.MethodDelete, "/gophers/1", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "cannot delete gopher")
}
