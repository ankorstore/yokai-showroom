package gopher_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ankorstore/yokai-showroom/http-demo/db/sqlc"
	"github.com/ankorstore/yokai-showroom/http-demo/internal"
	"github.com/ankorstore/yokai/fxsql"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestListGophersHandlerSuccess(t *testing.T) {
	var httpServer *echo.Echo
	var metricsRegistry *prometheus.Registry

	internal.RunTest(
		t,
		fxsql.RunFxSQLSeeds("gophers"),
		fx.Populate(&httpServer, &metricsRegistry),
	)

	// [GET] /gophers response assertion
	req := httptest.NewRequest(http.MethodGet, "/gophers", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var gophers []*sqlc.Gopher
	err := json.Unmarshal(rec.Body.Bytes(), &gophers)
	assert.NoError(t, err)

	assert.Len(t, gophers, 3)
	assert.Equal(t, gophers[0].Name, "alice")
	assert.Equal(t, gophers[0].Job.String, "architect")
	assert.Equal(t, gophers[1].Name, "bob")
	assert.Equal(t, gophers[1].Job.String, "builder")
	assert.Equal(t, gophers[2].Name, "carl")
	assert.Equal(t, gophers[2].Job.String, "carpenter")

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
}

func TestListGophersHandlerInternalServerErrorOnMissingTable(t *testing.T) {
	var db *sql.DB
	var httpServer *echo.Echo
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(t, fx.Populate(&db, &httpServer, &logBuffer, &traceExporter))

	// drop table for failure
	_, err := db.Exec("DROP TABLE gophers")
	assert.NoError(t, err)

	// [GET] /gophers response assertion
	req := httptest.NewRequest(http.MethodGet, "/gophers", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "cannot list gophers")
}
