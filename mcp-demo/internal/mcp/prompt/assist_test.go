package prompt_test

import (
	"context"
	"strings"
	"testing"

	"github.com/ankorstore/yokai-showroom/mcp-demo/internal"
	"github.com/ankorstore/yokai/fxmcpserver/fxmcpservertest"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/fx"
)

func TestJobAssistPrompt(t *testing.T) {
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

		internal.RunTest(t, fx.Populate(&testServer, &logBuffer, &traceExporter, &metricsRegistry))

		return testServer, logBuffer, traceExporter, metricsRegistry
	}

	t.Run("can get jobs assist prompt", func(t *testing.T) {
		testServer, logBuffer, traceExporter, metricsRegistry := runTest(t)

		defer testServer.Close()

		testClient, err := testServer.StartClient(context.Background())
		assert.NoError(t, err)

		defer testClient.Close()

		req := mcp.GetPromptRequest{}
		req.Params.Name = "jobs-assist"

		res, err := testClient.GetPrompt(context.Background(), req)
		assert.NoError(t, err)

		assert.Equal(t, "gophers jobs assistance", res.Description)

		assert.Equal(t, mcp.RoleAssistant, res.Messages[0].Role)
		assert.Equal(t, "help by providing information about the gopher jobs", res.Messages[0].Content.(mcp.TextContent).Text)

		embeddedResourceContents := res.Messages[1].Content.(mcp.EmbeddedResource).Resource.(mcp.TextResourceContents)
		assert.Equal(t, mcp.RoleAssistant, res.Messages[1].Role)
		assert.Equal(t, "jobs://jobs", embeddedResourceContents.URI)
		assert.Equal(t, "application/json", embeddedResourceContents.MIMEType)
		assert.Equal(t, `["backend","frontend"]`, embeddedResourceContents.Text)

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":        "info",
			"mcpMethod":    "prompts/get",
			"mcpPrompt":    "jobs-assist",
			"mcpTransport": "sse",
			"message":      "MCP request success",
		})

		tracetest.AssertHasTraceSpan(
			t,
			traceExporter,
			"MCP prompts/get jobs-assist",
			attribute.String("mcp.method", "prompts/get"),
			attribute.String("mcp.prompt", "jobs-assist"),
			attribute.String("mcp.transport", "sse"),
		)

		expectedMetric := `
        # HELP mcp_server_requests_total Number of processed MCP requests
        # TYPE mcp_server_requests_total counter
		mcp_server_requests_total{method="initialize",status="success",target=""} 1
        mcp_server_requests_total{method="prompts/get",status="success",target="jobs-assist"} 1
    `

		err = testutil.GatherAndCompare(
			metricsRegistry,
			strings.NewReader(expectedMetric),
			"mcp_server_requests_total",
		)
		assert.NoError(t, err)
	})
}
