package worker

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/prometheus/client_golang/prometheus"
)

// SubscribeCounter is a metrics counter for received messages.
var SubscribeCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "messages_received_total",
	Help: "Total number of received messages",
})

// SubscribeWorker is a worker to run pub/sub subscribers.
type SubscribeWorker struct {
	config *config.Config
	client *pubsub.Client
}

// NewSubscribeWorker returns a new SubscribeWorker.
func NewSubscribeWorker(config *config.Config, client *pubsub.Client) *SubscribeWorker {
	return &SubscribeWorker{
		config: config,
		client: client,
	}
}

// Name returns the worker name.
func (w *SubscribeWorker) Name() string {
	return "subscribe-worker"
}

// Run executes the worker.
func (w *SubscribeWorker) Run(ctx context.Context) error {
	subscription := w.client.Subscription(w.config.GetString("modules.pubsub.subscription"))

	return subscription.Receive(ctx, func(c context.Context, msg *pubsub.Message) {
		log.CtxLogger(c).Info().Msgf(
			"received message: id=%v, data=%v",
			msg.ID,
			string(msg.Data),
		)

		msg.Ack()

		SubscribeCounter.Inc()
	})
}
