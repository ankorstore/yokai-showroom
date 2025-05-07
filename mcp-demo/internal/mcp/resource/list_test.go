package resource_test

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

func TestJobsListResource(t *testing.T) {
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

	t.Run("can list jobs", func(t *testing.T) {
		testServer, logBuffer, traceExporter, metricsRegistry := runTest(t)

		defer testServer.Close()

		testClient, err := testServer.StartClient(context.Background())
		assert.NoError(t, err)

		defer testClient.Close()

		req := mcp.ReadResourceRequest{}
		req.Params.URI = "jobs://jobs"

		res, err := testClient.ReadResource(context.Background(), req)
		assert.NoError(t, err)

		resContents, ok := res.Contents[0].(mcp.TextResourceContents)
		assert.True(t, ok)

		assert.Equal(t, "jobs://jobs", resContents.URI)
		assert.Equal(t, "application/json", resContents.MIMEType)
		assert.Equal(t, `["backend","frontend"]`, resContents.Text)

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":          "info",
			"mcpMethod":      "resources/read",
			"mcpResourceURI": "jobs://jobs",
			"mcpTransport":   "sse",
			"message":        "MCP request success",
		})

		tracetest.AssertHasTraceSpan(
			t,
			traceExporter,
			"MCP resources/read jobs://jobs",
			attribute.String("mcp.method", "resources/read"),
			attribute.String("mcp.resourceURI", "jobs://jobs"),
			attribute.String("mcp.transport", "sse"),
		)

		expectedMetric := `
        # HELP mcp_server_requests_total Number of processed MCP requests
        # TYPE mcp_server_requests_total counter
		mcp_server_requests_total{method="initialize",status="success",target=""} 1
        mcp_server_requests_total{method="resources/read",status="success",target="jobs://jobs"} 1
    `

		err = testutil.GatherAndCompare(
			metricsRegistry,
			strings.NewReader(expectedMetric),
			"mcp_server_requests_total",
		)
		assert.NoError(t, err)
	})
}
