package repository_test

import (
	"context"
	"database/sql"
	"github.com/ankorstore/yokai-showroom/http-demo/internal"
	"github.com/ankorstore/yokai-showroom/http-demo/internal/repository"
	"github.com/ankorstore/yokai/fxsql"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"testing"
)

func TestFind(t *testing.T) {
	var db *sql.DB
	var gopherRepository *repository.GopherRepository
	var logger *log.Logger
	var tracerProvider oteltrace.TracerProvider
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	t.Run("should succeed", func(t *testing.T) {
		// run test
		internal.RunTest(
			t,
			fxsql.RunFxSQLSeeds(),
			fx.Populate(&db, &gopherRepository, &logger, &tracerProvider, &logBuffer, &traceExporter),
		)

		// context preparation
		ctx := logger.WithContext(trace.WithContext(context.Background(), tracerProvider))

		// result assertion
		gopher, err := gopherRepository.Find(ctx, 1)
		assert.NoError(t, err)

		assert.Equal(t, "alice", gopher.Name)
		assert.Equal(t, "frontend", gopher.Job.String)

		// log assertion
		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":     "debug",
			"system":    "sqlite",
			"operation": "connection:query-context",
			"query":     "SELECT id, name, job FROM gophers WHERE id = ? LIMIT 1",
			"arguments": "[map[Name: Ordinal:1 Value:1]]",
			"message":   "sql logger",
		})

		// trace assertion
		tracetest.AssertHasTraceSpan(
			t,
			traceExporter,
			"SQL connection:query-context",
			semconv.DBSystemKey.String("sqlite"),
			attribute.String("db.statement", "SELECT id, name, job FROM gophers WHERE id = ? LIMIT 1"),
			attribute.String("db.statement.arguments", "[{Name: Ordinal:1 Value:1}]"),
		)
	})

	t.Run("should fail", func(t *testing.T) {
		// run test
		internal.RunTest(t, fx.Populate(&db, &gopherRepository))

		// close db
		err := db.Close()
		assert.NoError(t, err)

		// result assertion
		_, err = gopherRepository.Find(context.Background(), 1)
		assert.Error(t, err)
	})
}

func TestFindAll(t *testing.T) {
	var db *sql.DB
	var gopherRepository *repository.GopherRepository
	var logger *log.Logger
	var tracerProvider oteltrace.TracerProvider
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	t.Run("should succeed", func(t *testing.T) {
		// run test
		internal.RunTest(
			t,
			fxsql.RunFxSQLSeeds(),
			fx.Populate(&db, &gopherRepository, &logger, &tracerProvider, &logBuffer, &traceExporter),
		)

		// context preparation
		ctx := logger.WithContext(trace.WithContext(context.Background(), tracerProvider))

		// result assertion
		gophers, err := gopherRepository.FindAll(ctx, repository.GopherRepositoryFindAllParams{
			Name: sql.NullString{String: "alice", Valid: true},
			Job:  sql.NullString{String: "frontend", Valid: true},
		})
		assert.NoError(t, err)

		assert.Len(t, gophers, 1)
		assert.Equal(t, gophers[0].Name, "alice")
		assert.Equal(t, gophers[0].Job.String, "frontend")

		// log assertion
		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":     "debug",
			"system":    "sqlite",
			"operation": "connection:query-context",
			"query":     "SELECT id, name, job FROM gophers WHERE name = ? AND job = ? ORDER BY id",
			"arguments": "[map[Name: Ordinal:1 Value:alice] map[Name: Ordinal:2 Value:frontend]]",
			"message":   "sql logger",
		})

		// trace assertion
		tracetest.AssertHasTraceSpan(
			t,
			traceExporter,
			"SQL connection:query-context",
			semconv.DBSystemKey.String("sqlite"),
			attribute.String("db.statement", "SELECT id, name, job FROM gophers WHERE name = ? AND job = ? ORDER BY id"),
			attribute.String("db.statement.arguments", "[{Name: Ordinal:1 Value:alice} {Name: Ordinal:2 Value:frontend}]"),
		)
	})

	t.Run("should fail", func(t *testing.T) {
		// run test
		internal.RunTest(t, fx.Populate(&db, &gopherRepository))

		// close db
		err := db.Close()
		assert.NoError(t, err)

		// result assertion
		_, err = gopherRepository.FindAll(context.Background(), repository.GopherRepositoryFindAllParams{})
		assert.Error(t, err)
	})
}

func TestCreate(t *testing.T) {
	var db *sql.DB
	var gopherRepository *repository.GopherRepository
	var logger *log.Logger
	var tracerProvider oteltrace.TracerProvider
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	t.Run("should succeed", func(t *testing.T) {
		// run test
		internal.RunTest(t, fx.Populate(&db, &gopherRepository, &logger, &tracerProvider, &logBuffer, &traceExporter))

		// context preparation
		ctx := logger.WithContext(trace.WithContext(context.Background(), tracerProvider))

		// result assertion
		gopherId, err := gopherRepository.Create(ctx, repository.GopherRepositoryCreateParams{
			Name: "test name",
			Job:  sql.NullString{String: "test job", Valid: true},
		})
		assert.NoError(t, err)

		assert.Equal(t, 1, gopherId)

		// log assertion
		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":     "debug",
			"system":    "sqlite",
			"operation": "connection:exec-context",
			"query":     "INSERT INTO gophers (name,job) VALUES (?,?)",
			"arguments": "[map[Name: Ordinal:1 Value:test name] map[Name: Ordinal:2 Value:test job]]",
			"message":   "sql logger",
		})

		// trace assertion
		tracetest.AssertHasTraceSpan(
			t,
			traceExporter,
			"SQL connection:exec-context",
			semconv.DBSystemKey.String("sqlite"),
			attribute.String("db.statement", "INSERT INTO gophers (name,job) VALUES (?,?)"),
			attribute.String("db.statement.arguments", "[{Name: Ordinal:1 Value:test name} {Name: Ordinal:2 Value:test job}]"),
		)
	})

	t.Run("should fail", func(t *testing.T) {
		// run test
		internal.RunTest(t, fx.Populate(&db, &gopherRepository))

		// close db
		err := db.Close()
		assert.NoError(t, err)

		// result assertion
		_, err = gopherRepository.Create(context.Background(), repository.GopherRepositoryCreateParams{})
		assert.Error(t, err)
	})
}

func TestDelete(t *testing.T) {
	var db *sql.DB
	var gopherRepository *repository.GopherRepository
	var logger *log.Logger
	var tracerProvider oteltrace.TracerProvider
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	t.Run("should succeed", func(t *testing.T) {
		// run test
		internal.RunTest(
			t,
			fxsql.RunFxSQLSeeds(),
			fx.Populate(&db, &gopherRepository, &logger, &tracerProvider, &logBuffer, &traceExporter),
		)

		// context preparation
		ctx := logger.WithContext(trace.WithContext(context.Background(), tracerProvider))

		// result assertion
		err := gopherRepository.Delete(ctx, 1)
		assert.NoError(t, err)

		// log assertion
		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":     "debug",
			"system":    "sqlite",
			"operation": "connection:exec-context",
			"query":     "DELETE FROM gophers WHERE id = ?",
			"arguments": "[map[Name: Ordinal:1 Value:1]]",
			"message":   "sql logger",
		})

		// trace assertion
		tracetest.AssertHasTraceSpan(
			t,
			traceExporter,
			"SQL connection:exec-context",
			semconv.DBSystemKey.String("sqlite"),
			attribute.String("db.statement", "DELETE FROM gophers WHERE id = ?"),
			attribute.String("db.statement.arguments", "[{Name: Ordinal:1 Value:1}]"),
		)
	})

	t.Run("should fail", func(t *testing.T) {
		// run test
		internal.RunTest(t, fx.Populate(&db, &gopherRepository))

		// close db
		err := db.Close()
		assert.NoError(t, err)

		// result assertion
		err = gopherRepository.Delete(context.Background(), 1)
		assert.Error(t, err)
	})
}
