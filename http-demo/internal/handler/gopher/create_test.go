package gopher_test

import (
	"bytes"
	"encoding/json"
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
	"gorm.io/gorm"
)

func TestCreateGopherHandlerSuccess(t *testing.T) {
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(t, fx.Populate(&httpServer, &logBuffer, &traceExporter))

	// [POST] /gophers response assertion
	data := `{"name": "gopher 1", "job": "job 1"}`
	req := httptest.NewRequest(http.MethodPost, "/gophers", bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var gopher *model.Gopher
	err := json.Unmarshal(rec.Body.Bytes(), &gopher)
	assert.NoError(t, err)

	assert.Equal(t, gopher.Name, "gopher 1")
	assert.Equal(t, gopher.Job, "job 1")

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
	assert.Contains(t, rec.Body.String(), "cannot bind gopher")
}

func TestCreateGopherHandlerInternalServerErrorOnMissingTable(t *testing.T) {
	var db *gorm.DB
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(t, fx.Populate(&db, &httpServer, &logBuffer, &traceExporter))

	// drop table for failure
	err := db.Migrator().DropTable(&model.Gopher{})
	assert.NoError(t, err)

	// [POST] /gophers response assertion
	req := httptest.NewRequest(http.MethodPost, "/gophers", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "cannot create gopher")
}
