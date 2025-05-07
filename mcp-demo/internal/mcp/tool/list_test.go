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

func TestListGophersTool(t *testing.T) {
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

	t.Run("can list all gophers", func(t *testing.T) {
		testServer, logBuffer, traceExporter, metricsRegistry := runTest(t)

		defer testServer.Close()

		testClient, err := testServer.StartClient(context.Background())
		assert.NoError(t, err)

		defer testClient.Close()

		req := mcp.CallToolRequest{}
		req.Params.Name = "list-gophers"

		res, err := testClient.CallTool(context.Background(), req)
		assert.NoError(t, err)
		assert.False(t, res.IsError)

		resContents, ok := res.Content[0].(mcp.TextContent)
		assert.True(t, ok)
		assert.Equal(
			t,
			`[{"id":1,"name":"alice","job":"frontend"},{"id":2,"name":"bob","job":"backend"},{"id":3,"name":"carl","job":"backend"},{"id":4,"name":"dan","job":"frontend"},{"id":5,"name":"elvis","job":"backend"}]`,
			resContents.Text,
		)

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":        "info",
			"mcpMethod":    "tools/call",
			"mcpTool":      "list-gophers",
			"mcpTransport": "sse",
			"message":      "MCP request success",
		})

		tracetest.AssertHasTraceSpan(
			t,
			traceExporter,
			"MCP tools/call list-gophers",
			attribute.String("mcp.method", "tools/call"),
			attribute.String("mcp.tool", "list-gophers"),
			attribute.String("mcp.transport", "sse"),
		)

		expectedMetric := `
        # HELP mcp_server_requests_total Number of processed MCP requests
        # TYPE mcp_server_requests_total counter
		mcp_server_requests_total{method="initialize",status="success",target=""} 1
        mcp_server_requests_total{method="tools/call",status="success",target="list-gophers"} 1
    `

		err = testutil.GatherAndCompare(
			metricsRegistry,
			strings.NewReader(expectedMetric),
			"mcp_server_requests_total",
		)
		assert.NoError(t, err)
	})

	t.Run("can list only frontend gophers", func(t *testing.T) {
		testServer, logBuffer, traceExporter, metricsRegistry := runTest(t)

		defer testServer.Close()

		testClient, err := testServer.StartClient(context.Background())
		assert.NoError(t, err)

		defer testClient.Close()

		req := mcp.CallToolRequest{}
		req.Params.Name = "list-gophers"
		req.Params.Arguments = map[string]any{
			"job": "frontend",
		}

		res, err := testClient.CallTool(context.Background(), req)
		assert.NoError(t, err)
		assert.False(t, res.IsError)

		resContents, ok := res.Content[0].(mcp.TextContent)
		assert.True(t, ok)
		assert.Equal(
			t,
			`[{"id":1,"name":"alice","job":"frontend"},{"id":4,"name":"dan","job":"frontend"}]`,
			resContents.Text,
		)

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":        "info",
			"mcpMethod":    "tools/call",
			"mcpTool":      "list-gophers",
			"mcpTransport": "sse",
			"message":      "MCP request success",
		})

		tracetest.AssertHasTraceSpan(
			t,
			traceExporter,
			"MCP tools/call list-gophers",
			attribute.String("mcp.method", "tools/call"),
			attribute.String("mcp.tool", "list-gophers"),
			attribute.String("mcp.transport", "sse"),
		)

		expectedMetric := `
        # HELP mcp_server_requests_total Number of processed MCP requests
        # TYPE mcp_server_requests_total counter
		mcp_server_requests_total{method="initialize",status="success",target=""} 1
        mcp_server_requests_total{method="tools/call",status="success",target="list-gophers"} 1
    `

		err = testutil.GatherAndCompare(
			metricsRegistry,
			strings.NewReader(expectedMetric),
			"mcp_server_requests_total",
		)
		assert.NoError(t, err)
	})
}
