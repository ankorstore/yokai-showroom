package internal

import (
	"github.com/ankorstore/yokai-showroom/worker-demo/internal/service"
	"github.com/ankorstore/yokai-showroom/worker-demo/internal/worker"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxworker"
	"go.uber.org/fx"
)

func ProvideServices() fx.Option {
	return fx.Options(
		// annotated publisher service
		fx.Provide(
			fx.Annotate(
				service.NewDefaultPublisher,
				fx.As(new(service.Publisher)),
			),
		),
		// annotated subscriber service
		fx.Provide(
			fx.Annotate(
				service.NewDefaultSubscriber,
				fx.As(new(service.Subscriber)),
			),
		),
		// subscriber worker
		fxworker.AsWorker(worker.NewSubscribeWorker),
		// metrics
		fxmetrics.AsMetricsCollectors(
			service.PublishCounter,
			service.SubscribeCounter,
		),
	)
}
