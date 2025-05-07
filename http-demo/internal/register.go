package internal

import (
	"github.com/ankorstore/yokai-showroom/http-demo/internal/domain"
	"github.com/ankorstore/yokai/fxhealthcheck"
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
		// metrics
		fxmetrics.AsMetricsCollector(domain.GopherServiceCounter),
		// probes
		fxhealthcheck.AsCheckerProbe(healthcheck.NewSQLProbe),
	)
}
