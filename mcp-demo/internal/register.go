package internal

import (
	"github.com/ankorstore/yokai-showroom/mcp-demo/internal/domain"
	"github.com/ankorstore/yokai-showroom/mcp-demo/internal/mcp/prompt"
	"github.com/ankorstore/yokai-showroom/mcp-demo/internal/mcp/resource"
	"github.com/ankorstore/yokai-showroom/mcp-demo/internal/mcp/tool"
	"github.com/ankorstore/yokai/fxhealthcheck"
	"github.com/ankorstore/yokai/fxmcpserver"
	"github.com/ankorstore/yokai/fxmcpserver/server"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/sql/healthcheck"
	"go.uber.org/fx"
)

// Register is used to register the application dependencies.
func Register() fx.Option {
	return fx.Options(
		// domain
		fx.Provide(
			domain.NewGopherRepository,
			domain.NewGopherService,
		),
		// mcp
		fxmcpserver.AsMCPServerPrompt(prompt.NewJobsAssistPrompt),
		fxmcpserver.AsMCPServerResource(resource.NewJobsListResource),
		fxmcpserver.AsMCPServerResourceTemplate(resource.NewJobDetailResource),
		fxmcpserver.AsMCPServerTools(
			tool.NewListGophersTool,
			tool.NewGetGopherTool,
			tool.NewCreateGopherTool,
			tool.NewDeleteGopherTool,
		),
		// metrics
		fxmetrics.AsMetricsCollector(domain.GopherServiceCounter),
		// probes
		fxhealthcheck.AsCheckerProbe(healthcheck.NewSQLProbe),
		fxhealthcheck.AsCheckerProbe(server.NewMCPServerProbe),
	)
}
