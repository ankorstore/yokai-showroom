package internal

import (
	"github.com/ankorstore/yokai-showroom/http-demo/internal/handler"
	"github.com/ankorstore/yokai-showroom/http-demo/internal/repository"
	"github.com/ankorstore/yokai-showroom/http-demo/internal/service"
	"github.com/ankorstore/yokai/fxhealthcheck"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/orm/healthcheck"
	"go.uber.org/fx"
)

func ProvideServices() fx.Option {
	return fx.Options(
		// dashboard metrics
		fxmetrics.AsMetricsCollector(handler.DashboardHistogram),
		// orm healthcheck probe
		fxhealthcheck.AsCheckerProbe(healthcheck.NewOrmProbe),
		// services
		fx.Provide(
			// gophers repository
			repository.NewGopherRepository,
			// gophers service
			service.NewGopherService,
		),
	)
}
