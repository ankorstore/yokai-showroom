package gopher_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankorstore/yokai-showroom/http-demo/internal"
	"github.com/ankorstore/yokai-showroom/http-demo/internal/model"
	"github.com/ankorstore/yokai-showroom/http-demo/internal/repository"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func TestUpdateGopherHandlerSuccess(t *testing.T) {
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter
	var repository *repository.GopherRepository

	internal.RunTest(t, fx.Populate(&httpServer, &logBuffer, &traceExporter, &repository))

	// populate database
	err := repository.Create(context.Background(), &model.Gopher{
		Name: "gopher 1",
		Job:  "job 1",
	})
	assert.NoError(t, err)

	// [PATCH] /gophers/1 response assertion
	data := `{"name": "new gopher 1", "job": "new job 1"}`
	req := httptest.NewRequest(http.MethodPatch, "/gophers/1", bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var gopher *model.Gopher
	err = json.Unmarshal(rec.Body.Bytes(), &gopher)
	assert.NoError(t, err)

	assert.Equal(t, gopher.Name, "new gopher 1")
	assert.Equal(t, gopher.Job, "new job 1")

	// log assertion
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "called update gopher",
	})

	// trace assertion
	tracetest.AssertHasTraceSpan(t, traceExporter, "update gopher service")
}

func TestUpdateGopherHandlerBadRequestErrorOnInvalidId(t *testing.T) {
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(t, fx.Populate(&httpServer, &logBuffer, &traceExporter))

	// [PATCH] /gophers/1 response assertion
	data := `invalid`
	req := httptest.NewRequest(http.MethodPatch, "/gophers/1", bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "cannot bind gopher")
}

func TestUpdateGopherHandlerBadRequestErrorOnInvalidData(t *testing.T) {
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(t, fx.Populate(&httpServer, &logBuffer, &traceExporter))

	// [PATCH] /gophers/invalid response assertion
	req := httptest.NewRequest(http.MethodPatch, "/gophers/invalid", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid gopher id")
}

func TestUpdateGopherHandlerNotFoundErrorOnMissingId(t *testing.T) {
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(t, fx.Populate(&httpServer, &logBuffer, &traceExporter))

	// [PATCH] /gophers/1 response assertion
	req := httptest.NewRequest(http.MethodPatch, "/gophers/1", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "cannot get gopher with id 1")
}

func TestUpdateGopherHandlerInternalServerErrorOnMissingTable(t *testing.T) {
	var db *gorm.DB
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(t, fx.Populate(&db, &httpServer, &logBuffer, &traceExporter))

	// drop table for failure
	err := db.Migrator().DropTable(&model.Gopher{})
	assert.NoError(t, err)

	// [PATCH] /gophers/1 response assertion
	req := httptest.NewRequest(http.MethodPatch, "/gophers/1", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "cannot update gopher")
}
