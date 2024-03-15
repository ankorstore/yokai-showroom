package internal

import (
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/healthcheck"
	"github.com/ankorstore/yokai-showroom/worker-demo/internal/worker"
	"github.com/ankorstore/yokai/fxhealthcheck"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxworker"
	"go.uber.org/fx"
)

func ProvideServices() fx.Option {
	return fx.Options(
		// pub/sub subscriber worker
		fxworker.AsWorker(worker.NewSubscribeWorker),
		// pub/sub subscriber metrics
		fxmetrics.AsMetricsCollector(worker.SubscribeCounter),
		// pub/sub subscription health check
		fxhealthcheck.AsCheckerProbe(healthcheck.NewGcpPubSubSubscriptionsProbe),
	)
}
