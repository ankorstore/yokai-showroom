package service_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai-showroom/http-demo/internal"
	"github.com/ankorstore/yokai-showroom/http-demo/internal/model"
	"github.com/ankorstore/yokai-showroom/http-demo/internal/repository"
	"github.com/ankorstore/yokai-showroom/http-demo/internal/service"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestList(t *testing.T) {
	var repo *repository.GopherRepository
	var svc *service.GopherService
	var logger *log.Logger
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(t, fx.Populate(&repo, &svc, &logger, &logBuffer, &traceExporter))

	ctx := logger.WithContext(context.Background())

	gopher := &model.Gopher{
		Name: "gopher 1",
		Job:  "job 1",
	}

	err := repo.Create(ctx, gopher)
	assert.NoError(t, err)

	foundGophers, err := svc.List(ctx)
	assert.NoError(t, err)
	assert.Len(t, foundGophers, 1)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "called list gophers",
	})

	tracetest.AssertHasTraceSpan(t, traceExporter, "list gophers service")
}

func TestGet(t *testing.T) {
	var repo *repository.GopherRepository
	var svc *service.GopherService
	var logger *log.Logger
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(t, fx.Populate(&repo, &svc, &logger, &logBuffer, &traceExporter))

	ctx := logger.WithContext(context.Background())

	gopher := &model.Gopher{
		Name: "gopher 1",
		Job:  "job 1",
	}

	err := repo.Create(ctx, gopher)
	assert.NoError(t, err)

	foundGopher, err := svc.Get(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "gopher 1", foundGopher.Name)
	assert.Equal(t, "job 1", foundGopher.Job)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "called get gopher",
	})

	tracetest.AssertHasTraceSpan(t, traceExporter, "get gopher service")
}

func TestCreate(t *testing.T) {
	var svc *service.GopherService
	var logger *log.Logger
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(t, fx.Populate(&svc, &logger, &logBuffer, &traceExporter))

	ctx := logger.WithContext(context.Background())

	gopher := &model.Gopher{
		Name: "gopher 1",
		Job:  "job 1",
	}

	err := svc.Create(ctx, gopher)
	assert.NoError(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "called create gopher",
	})

	tracetest.AssertHasTraceSpan(t, traceExporter, "create gopher service")
}

func TestDelete(t *testing.T) {
	var repo *repository.GopherRepository
	var svc *service.GopherService
	var logger *log.Logger
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(t, fx.Populate(&repo, &svc, &logger, &logBuffer, &traceExporter))

	ctx := logger.WithContext(context.Background())

	gopher := &model.Gopher{
		Name: "gopher 1",
		Job:  "job 1",
	}

	err := repo.Create(ctx, gopher)
	assert.NoError(t, err)

	err = svc.Delete(ctx, 1)
	assert.NoError(t, err)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "called delete gopher",
	})

	tracetest.AssertHasTraceSpan(t, traceExporter, "delete gopher service")
}

func TestUpdate(t *testing.T) {
	var repo *repository.GopherRepository
	var svc *service.GopherService
	var logger *log.Logger
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(t, fx.Populate(&repo, &svc, &logger, &logBuffer, &traceExporter))

	ctx := logger.WithContext(context.Background())

	gopher := &model.Gopher{
		Name: "gopher 1",
		Job:  "job 1",
	}

	err := repo.Create(ctx, gopher)
	assert.NoError(t, err)

	updatedGopher, err := svc.Update(ctx, 1, &model.Gopher{
		Name: "new gopher 1",
		Job:  "new job 1",
	})
	assert.NoError(t, err)
	assert.Equal(t, "new gopher 1", updatedGopher.Name)
	assert.Equal(t, "new job 1", updatedGopher.Job)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "called update gopher",
	})

	tracetest.AssertHasTraceSpan(t, traceExporter, "update gopher service")
}
