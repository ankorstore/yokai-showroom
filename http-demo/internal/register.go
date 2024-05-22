package internal

import (
	"github.com/ankorstore/yokai-showroom/http-demo/internal/repository"
	"github.com/ankorstore/yokai-showroom/http-demo/internal/service"
	"github.com/ankorstore/yokai/fxhealthcheck"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/sql/healthcheck"
	"go.uber.org/fx"
)

// Register is used to register the application dependencies.
func Register() fx.Option {
	return fx.Options(
		// services
		fx.Provide(
			repository.NewGopherRepository,
			service.NewGopherService,
		),
		// metrics
		fxmetrics.AsMetricsCollector(service.GopherServiceCounter),
		// probes
		fxhealthcheck.AsCheckerProbe(healthcheck.NewSQLProbe),
	)
}
