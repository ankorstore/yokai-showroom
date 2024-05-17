package internal

import (
	"github.com/ankorstore/yokai-showroom/http-demo/internal/handler"
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
		// metrics
		fxmetrics.AsMetricsCollectors(
			handler.DashboardHistogram,
			service.GopherServiceCounter,
		),
		// healthcheck probe
		fxhealthcheck.AsCheckerProbe(healthcheck.NewSQLProbe),
		// services
		fx.Provide(
			// gophers repository
			repository.NewGopherRepository,
			// gophers service
			service.NewGopherService,
		),
	)
}
