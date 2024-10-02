package internal

import (
	pubsubhc "github.com/ankorstore/yokai-contrib/fxgcppubsub/healthcheck"
	"github.com/ankorstore/yokai-showroom/worker-demo/internal/worker"
	"github.com/ankorstore/yokai/fxhealthcheck"
	"github.com/ankorstore/yokai/fxmetrics"
	"github.com/ankorstore/yokai/fxworker"
	workerhc "github.com/ankorstore/yokai/worker/healthcheck"
	"go.uber.org/fx"
)

// Register is used to register the application dependencies.
func Register() fx.Option {
	return fx.Options(
		// pub/sub subscriber worker
		fxworker.AsWorker(worker.NewSubscribeWorker),
		// pub/sub subscriber metrics
		fxmetrics.AsMetricsCollector(worker.SubscribeCounter),
		// pub/sub subscription health check
		fxhealthcheck.AsCheckerProbe(pubsubhc.NewGcpPubSubSubscriptionsProbe),
		// worker pool health check
		fxhealthcheck.AsCheckerProbe(workerhc.NewWorkerProbe),
	)
}
