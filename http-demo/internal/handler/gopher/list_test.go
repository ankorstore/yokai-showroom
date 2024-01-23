package gopher_test

import (
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

func TestListGophersHandlerSuccess(t *testing.T) {
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

	err = repository.Create(context.Background(), &model.Gopher{
		Name: "gopher 2",
		Job:  "job 2",
	})
	assert.NoError(t, err)

	// [GET] /gophers response assertion
	req := httptest.NewRequest(http.MethodGet, "/gophers", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var gophers []*model.Gopher
	err = json.Unmarshal(rec.Body.Bytes(), &gophers)
	assert.NoError(t, err)

	assert.Len(t, gophers, 2)
	assert.Equal(t, gophers[0].Name, "gopher 1")
	assert.Equal(t, gophers[0].Job, "job 1")
	assert.Equal(t, gophers[1].Name, "gopher 2")
	assert.Equal(t, gophers[1].Job, "job 2")

	// log assertion
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "called list gophers",
	})

	// trace assertion
	tracetest.AssertHasTraceSpan(t, traceExporter, "list gophers service")
}

func TestListGophersHandlerInternalServerErrorOnMissingTable(t *testing.T) {
	var db *gorm.DB
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(t, fx.Populate(&db, &httpServer, &logBuffer, &traceExporter))

	// drop table for failure
	err := db.Migrator().DropTable(&model.Gopher{})
	assert.NoError(t, err)

	// [GET] /gophers response assertion
	req := httptest.NewRequest(http.MethodGet, "/gophers", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "cannot list gophers")
}
