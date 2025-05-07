package domain_test

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/ankorstore/yokai-showroom/http-demo/internal"
	"github.com/ankorstore/yokai-showroom/http-demo/internal/domain"
	"github.com/ankorstore/yokai/fxsql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestGopherService_Get(t *testing.T) {
	var db *sql.DB
	var gopherService *domain.GopherService
	var metricsRegistry *prometheus.Registry

	t.Run("should succeed", func(t *testing.T) {
		// reset
		domain.GopherServiceCounter.Reset()

		// run test
		internal.RunTest(
			t,
			fxsql.RunFxSQLSeeds(),
			fx.Populate(&db, &gopherService, &metricsRegistry),
		)

		// result assertion
		gopher, err := gopherService.Get(context.Background(), 1)
		assert.NoError(t, err)

		assert.Equal(t, "alice", gopher.Name)
		assert.Equal(t, "frontend", gopher.Job.String)

		// metrics assertion
		expectedMetric := `
			# HELP gophers_service_operations_total Number of operations on the GopherService
			# TYPE gophers_service_operations_total counter
			gophers_service_operations_total{operation="get"} 1
		`

		err = testutil.GatherAndCompare(
			metricsRegistry,
			strings.NewReader(expectedMetric),
			"gophers_service_operations_total",
		)
		assert.NoError(t, err)
	})

	t.Run("should fail", func(t *testing.T) {
		// run test
		internal.RunTest(t, fx.Populate(&db, &gopherService))

		// close db
		err := db.Close()
		assert.NoError(t, err)

		// result assertion
		_, err = gopherService.Get(context.Background(), 1)
		assert.Error(t, err)
	})
}

func TestGopherService_List(t *testing.T) {
	var db *sql.DB
	var gopherService *domain.GopherService
	var metricsRegistry *prometheus.Registry

	t.Run("should succeed", func(t *testing.T) {
		// reset
		domain.GopherServiceCounter.Reset()

		// run test
		internal.RunTest(
			t,
			fxsql.RunFxSQLSeeds(),
			fx.Populate(&db, &gopherService, &metricsRegistry),
		)

		// result assertion
		gophers, err := gopherService.List(context.Background(), "alice", "frontend")
		assert.NoError(t, err)

		assert.Len(t, gophers, 1)
		assert.Equal(t, "alice", gophers[0].Name)
		assert.Equal(t, "frontend", gophers[0].Job.String)

		// metrics assertion
		expectedMetric := `
			# HELP gophers_service_operations_total Number of operations on the GopherService
			# TYPE gophers_service_operations_total counter
			gophers_service_operations_total{operation="list"} 1
		`

		err = testutil.GatherAndCompare(
			metricsRegistry,
			strings.NewReader(expectedMetric),
			"gophers_service_operations_total",
		)
		assert.NoError(t, err)
	})

	t.Run("should fail", func(t *testing.T) {
		// run test
		internal.RunTest(t, fx.Populate(&db, &gopherService))

		// close db
		err := db.Close()
		assert.NoError(t, err)

		// result assertion
		_, err = gopherService.List(context.Background(), "", "")
		assert.Error(t, err)
	})
}

func TestGopherService_Create(t *testing.T) {
	var db *sql.DB
	var gopherService *domain.GopherService
	var metricsRegistry *prometheus.Registry

	t.Run("should succeed", func(t *testing.T) {
		// reset
		domain.GopherServiceCounter.Reset()

		// run test
		internal.RunTest(
			t,
			fx.Populate(&db, &gopherService, &metricsRegistry),
		)

		// result assertion
		gopherId, err := gopherService.Create(context.Background(), "alice", "frontend")
		assert.NoError(t, err)

		assert.Equal(t, 1, gopherId)

		// metrics assertion
		expectedMetric := `
			# HELP gophers_service_operations_total Number of operations on the GopherService
			# TYPE gophers_service_operations_total counter
			gophers_service_operations_total{operation="create"} 1
		`

		err = testutil.GatherAndCompare(
			metricsRegistry,
			strings.NewReader(expectedMetric),
			"gophers_service_operations_total",
		)
		assert.NoError(t, err)
	})

	t.Run("should fail", func(t *testing.T) {
		// run test
		internal.RunTest(t, fx.Populate(&db, &gopherService))

		// close db
		err := db.Close()
		assert.NoError(t, err)

		// result assertion
		_, err = gopherService.Create(context.Background(), "", "")
		assert.Error(t, err)
	})
}

func TestGopherService_Delete(t *testing.T) {
	var db *sql.DB
	var gopherService *domain.GopherService
	var metricsRegistry *prometheus.Registry

	t.Run("should succeed", func(t *testing.T) {
		// reset
		domain.GopherServiceCounter.Reset()

		// run test
		internal.RunTest(
			t,
			fxsql.RunFxSQLSeeds(),
			fx.Populate(&db, &gopherService, &metricsRegistry),
		)

		// result assertion
		err := gopherService.Delete(context.Background(), 1)
		assert.NoError(t, err)

		// metrics assertion
		expectedMetric := `
			# HELP gophers_service_operations_total Number of operations on the GopherService
			# TYPE gophers_service_operations_total counter
			gophers_service_operations_total{operation="delete"} 1
		`

		err = testutil.GatherAndCompare(
			metricsRegistry,
			strings.NewReader(expectedMetric),
			"gophers_service_operations_total",
		)
		assert.NoError(t, err)
	})

	t.Run("should fail", func(t *testing.T) {
		// run test
		internal.RunTest(t, fx.Populate(&db, &gopherService))

		// close db
		err := db.Close()
		assert.NoError(t, err)

		// result assertion
		err = gopherService.Delete(context.Background(), 1)
		assert.Error(t, err)
	})
}
