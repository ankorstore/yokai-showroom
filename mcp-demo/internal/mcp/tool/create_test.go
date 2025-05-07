package tool_test

import (
	"context"
	"strings"
	"testing"

	"github.com/ankorstore/yokai-showroom/mcp-demo/internal"
	"github.com/ankorstore/yokai-showroom/mcp-demo/internal/domain"
	"github.com/ankorstore/yokai/fxmcpserver/fxmcpservertest"
	"github.com/ankorstore/yokai/fxsql"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/fx"
)

func TestCreateGopherTool(t *testing.T) {
	runTest := func(
		tb testing.TB,
	) (
		*fxmcpservertest.MCPSSETestServer,
		logtest.TestLogBuffer,
		tracetest.TestTraceExporter,
		*prometheus.Registry,
		*domain.GopherRepository,
	) {
		tb.Helper()

		var testServer *fxmcpservertest.MCPSSETestServer
		var logBuffer logtest.TestLogBuffer
		var traceExporter tracetest.TestTraceExporter
		var metricsRegistry *prometheus.Registry
		var gopherRepository *domain.GopherRepository

		internal.RunTest(
			t,
			fxsql.RunFxSQLSeeds(),
			fx.Populate(&testServer, &logBuffer, &traceExporter, &metricsRegistry, &gopherRepository),
		)

		return testServer, logBuffer, traceExporter, metricsRegistry, gopherRepository
	}

	t.Run("can create a new gopher", func(t *testing.T) {
		testServer, logBuffer, traceExporter, metricsRegistry, gopherRepository := runTest(t)

		defer testServer.Close()

		testClient, err := testServer.StartClient(context.Background())
		assert.NoError(t, err)

		defer testClient.Close()

		req := mcp.CallToolRequest{}
		req.Params.Name = "create-gopher"
		req.Params.Arguments = map[string]any{
			"name": "zoe",
			"job":  "backend",
		}

		res, err := testClient.CallTool(context.Background(), req)
		assert.NoError(t, err)
		assert.False(t, res.IsError)

		resContents, ok := res.Content[0].(mcp.TextContent)
		assert.True(t, ok)
		assert.Equal(t, "gopher was created with id 6", resContents.Text)

		gopher, err := gopherRepository.Find(context.Background(), 6)
		assert.NoError(t, err)

		assert.Equal(t, "zoe", gopher.Name)
		assert.Equal(t, "backend", gopher.Job)

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":        "info",
			"mcpMethod":    "tools/call",
			"mcpTool":      "create-gopher",
			"mcpTransport": "sse",
			"message":      "MCP request success",
		})

		tracetest.AssertHasTraceSpan(
			t,
			traceExporter,
			"MCP tools/call create-gopher",
			attribute.String("mcp.method", "tools/call"),
			attribute.String("mcp.tool", "create-gopher"),
			attribute.String("mcp.transport", "sse"),
		)

		expectedMetric := `
        # HELP mcp_server_requests_total Number of processed MCP requests
        # TYPE mcp_server_requests_total counter
		mcp_server_requests_total{method="initialize",status="success",target=""} 1
        mcp_server_requests_total{method="tools/call",status="success",target="create-gopher"} 1
    `

		err = testutil.GatherAndCompare(
			metricsRegistry,
			strings.NewReader(expectedMetric),
			"mcp_server_requests_total",
		)
		assert.NoError(t, err)
	})
}
