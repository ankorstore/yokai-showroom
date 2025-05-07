package tool_test

import (
	"context"
	"strings"
	"testing"

	"github.com/ankorstore/yokai-showroom/mcp-demo/internal"
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

func TestDeleteGopherTool(t *testing.T) {
	runTest := func(
		tb testing.TB,
	) (
		*fxmcpservertest.MCPSSETestServer,
		logtest.TestLogBuffer,
		tracetest.TestTraceExporter,
		*prometheus.Registry,
	) {
		tb.Helper()

		var testServer *fxmcpservertest.MCPSSETestServer
		var logBuffer logtest.TestLogBuffer
		var traceExporter tracetest.TestTraceExporter
		var metricsRegistry *prometheus.Registry

		internal.RunTest(
			t,
			fxsql.RunFxSQLSeeds(),
			fx.Populate(&testServer, &logBuffer, &traceExporter, &metricsRegistry),
		)

		return testServer, logBuffer, traceExporter, metricsRegistry
	}

	t.Run("can delete existing gopher", func(t *testing.T) {
		testServer, logBuffer, traceExporter, metricsRegistry := runTest(t)

		defer testServer.Close()

		testClient, err := testServer.StartClient(context.Background())
		assert.NoError(t, err)

		defer testClient.Close()

		req := mcp.CallToolRequest{}
		req.Params.Name = "delete-gopher"
		req.Params.Arguments = map[string]any{
			"id": "2",
		}

		res, err := testClient.CallTool(context.Background(), req)
		assert.NoError(t, err)
		assert.False(t, res.IsError)

		resContents, ok := res.Content[0].(mcp.TextContent)
		assert.True(t, ok)
		assert.Equal(t, "gopher with id 2 was deleted", resContents.Text)

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":        "info",
			"mcpMethod":    "tools/call",
			"mcpTool":      "delete-gopher",
			"mcpTransport": "sse",
			"message":      "MCP request success",
		})

		tracetest.AssertHasTraceSpan(
			t,
			traceExporter,
			"MCP tools/call delete-gopher",
			attribute.String("mcp.method", "tools/call"),
			attribute.String("mcp.tool", "delete-gopher"),
			attribute.String("mcp.transport", "sse"),
		)

		expectedMetric := `
        # HELP mcp_server_requests_total Number of processed MCP requests
        # TYPE mcp_server_requests_total counter
		mcp_server_requests_total{method="initialize",status="success",target=""} 1
        mcp_server_requests_total{method="tools/call",status="success",target="delete-gopher"} 1
    `

		err = testutil.GatherAndCompare(
			metricsRegistry,
			strings.NewReader(expectedMetric),
			"mcp_server_requests_total",
		)
		assert.NoError(t, err)
	})

	t.Run("cannot delete invalid gopher", func(t *testing.T) {
		testServer, logBuffer, traceExporter, metricsRegistry := runTest(t)

		defer testServer.Close()

		testClient, err := testServer.StartClient(context.Background())
		assert.NoError(t, err)

		defer testClient.Close()

		req := mcp.CallToolRequest{}
		req.Params.Name = "delete-gopher"
		req.Params.Arguments = map[string]any{
			"id": "99",
		}

		_, err = testClient.CallTool(context.Background(), req)
		assert.Error(t, err)
		assert.Equal(t, "cannot get gopher: sql: no rows in result set", err.Error())

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":        "error",
			"mcpError":     "request error: cannot get gopher: sql: no rows in result set",
			"mcpMethod":    "tools/call",
			"mcpTool":      "delete-gopher",
			"mcpTransport": "sse",
			"message":      "MCP request error",
		})

		tracetest.AssertHasTraceSpan(
			t,
			traceExporter,
			"MCP tools/call delete-gopher",
			attribute.String("mcp.error", "request error: cannot get gopher: sql: no rows in result set"),
			attribute.String("mcp.method", "tools/call"),
			attribute.String("mcp.tool", "delete-gopher"),
			attribute.String("mcp.transport", "sse"),
		)

		expectedMetric := `
        # HELP mcp_server_requests_total Number of processed MCP requests
        # TYPE mcp_server_requests_total counter
		mcp_server_requests_total{method="initialize",status="success",target=""} 1
        mcp_server_requests_total{method="tools/call",status="error",target="delete-gopher"} 1
    `

		err = testutil.GatherAndCompare(
			metricsRegistry,
			strings.NewReader(expectedMetric),
			"mcp_server_requests_total",
		)
		assert.NoError(t, err)
	})
}
