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

func TestJobDetailResource(t *testing.T) {
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

	t.Run("can get details of valid job", func(t *testing.T) {
		testServer, logBuffer, traceExporter, metricsRegistry := runTest(t)

		defer testServer.Close()

		testClient, err := testServer.StartClient(context.Background())
		assert.NoError(t, err)

		defer testClient.Close()

		req := mcp.ReadResourceRequest{}
		req.Params.URI = "jobs://jobs/frontend"

		res, err := testClient.ReadResource(context.Background(), req)
		assert.NoError(t, err)

		resContents, ok := res.Contents[0].(mcp.TextResourceContents)
		assert.True(t, ok)

		assert.Equal(t, "jobs://jobs/frontend", resContents.URI)
		assert.Equal(t, "text/plain", resContents.MIMEType)
		assert.Equal(t, "This is the job for gophers that love clean interfaces.", resContents.Text)

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":          "info",
			"mcpMethod":      "resources/read",
			"mcpResourceURI": "jobs://jobs/frontend",
			"mcpTransport":   "sse",
			"message":        "MCP request success",
		})

		tracetest.AssertHasTraceSpan(
			t,
			traceExporter,
			"MCP resources/read jobs://jobs/frontend",
			attribute.String("mcp.method", "resources/read"),
			attribute.String("mcp.resourceURI", "jobs://jobs/frontend"),
			attribute.String("mcp.transport", "sse"),
		)

		expectedMetric := `
        # HELP mcp_server_requests_total Number of processed MCP requests
        # TYPE mcp_server_requests_total counter
		mcp_server_requests_total{method="initialize",status="success",target=""} 1
        mcp_server_requests_total{method="resources/read",status="success",target="jobs://jobs/frontend"} 1
    `

		err = testutil.GatherAndCompare(
			metricsRegistry,
			strings.NewReader(expectedMetric),
			"mcp_server_requests_total",
		)
		assert.NoError(t, err)
	})

	t.Run("cannot get details of invalid job", func(t *testing.T) {
		testServer, logBuffer, traceExporter, metricsRegistry := runTest(t)

		defer testServer.Close()

		testClient, err := testServer.StartClient(context.Background())
		assert.NoError(t, err)

		defer testClient.Close()

		req := mcp.ReadResourceRequest{}
		req.Params.URI = "jobs://jobs/invalid"

		_, err = testClient.ReadResource(context.Background(), req)
		assert.Error(t, err)
		assert.Equal(t, "job invalid does not exist", err.Error())

		logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
			"level":          "error",
			"mcpError":       "request error: job invalid does not exist",
			"mcpMethod":      "resources/read",
			"mcpResourceURI": "jobs://jobs/invalid",
			"mcpTransport":   "sse",
			"message":        "MCP request error",
		})

		tracetest.AssertHasTraceSpan(
			t,
			traceExporter,
			"MCP resources/read jobs://jobs/invalid",
			attribute.String("mcp.error", "request error: job invalid does not exist"),
			attribute.String("mcp.method", "resources/read"),
			attribute.String("mcp.resourceURI", "jobs://jobs/invalid"),
			attribute.String("mcp.transport", "sse"),
		)

		expectedMetric := `
        # HELP mcp_server_requests_total Number of processed MCP requests
        # TYPE mcp_server_requests_total counter
		mcp_server_requests_total{method="initialize",status="success",target=""} 1
        mcp_server_requests_total{method="resources/read",status="error",target="jobs://jobs/invalid"} 1
    `

		err = testutil.GatherAndCompare(
			metricsRegistry,
			strings.NewReader(expectedMetric),
			"mcp_server_requests_total",
		)
		assert.NoError(t, err)
	})
}
