package worker

import (
	"context"
	"fmt"

	"github.com/ankorstore/yokai-contrib/fxgcppubsub"
	"github.com/ankorstore/yokai-contrib/fxgcppubsub/message"
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/worker"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// SubscribeCounter is a metrics counter for received messages.
var SubscribeCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "subscriber_messages_received_total",
	Help: "Total number of received messages",
})

// SubscribeWorker is a worker to run pub/sub subscribers.
type SubscribeWorker struct {
	config     *config.Config
	subscriber fxgcppubsub.Subscriber
}

// NewSubscribeWorker returns a new SubscribeWorker.
func NewSubscribeWorker(config *config.Config, subscriber fxgcppubsub.Subscriber) *SubscribeWorker {
	return &SubscribeWorker{
		config:     config,
		subscriber: subscriber,
	}
}

// Name returns the worker name.
func (w *SubscribeWorker) Name() string {
	return "subscribe-worker"
}

// Run executes the worker.
func (w *SubscribeWorker) Run(ctx context.Context) error {
	tracer := worker.CtxTracer(ctx)

	return w.subscriber.Subscribe(ctx, w.config.GetString("config.subscription.id"), func(fCtx context.Context, msg *message.Message) {
		data := string(msg.Data())

		fCtx, span := tracer.Start(
			fCtx,
			fmt.Sprintf("%s message", w.Name()),
			oteltrace.WithNewRoot(),
			oteltrace.WithSpanKind(oteltrace.SpanKindConsumer),
			oteltrace.WithAttributes(attribute.String("Message", data)),
		)
		defer span.End()

		log.CtxLogger(fCtx).Info().Msgf("received message: id=%v, data=%v", msg.ID(), data)

		msg.Ack()

		SubscribeCounter.Inc()
	})
}
