package internal

import (
	"github.com/ankorstore/yokai-showroom/worker-demo/internal/worker"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxworker"
	"go.uber.org/fx"
)

func ProvideServices() fx.Option {
	return fx.Options(
		// subscriber worker
		fxworker.AsWorker(worker.NewSubscribeWorker),
		// subscriber metrics
		fxmetrics.AsMetricsCollector(worker.SubscribeCounter),
	)
}
