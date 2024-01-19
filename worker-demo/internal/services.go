package internal

import (
	"github.com/ankorstore/yokai-showroom/worker-demo/internal/service"
	"github.com/ankorstore/yokai-showroom/worker-demo/internal/worker"
	"github.com/ankorstore/yokai/fxworker"
	"go.uber.org/fx"
)

func ProvideServices() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				service.NewPubSubPublisher,
				fx.As(new(service.Publisher)),
			),
		),
		fx.Provide(
			fx.Annotate(
				service.NewPubSubSubscriber,
				fx.As(new(service.Subscriber)),
			),
		),
		fxworker.AsWorker(worker.NewSubscribeWorker),
	)
}
