package worker

import (
	"context"

	"github.com/ankorstore/yokai-showroom/worker-demo/internal/service"
)

// SubscribeWorker is a worker to run pub/sub subscribers.
type SubscribeWorker struct {
	subscriber service.Subscriber
}

// NewSubscribeWorker returns a new SubscribeWorker.
func NewSubscribeWorker(subscriber service.Subscriber) *SubscribeWorker {
	return &SubscribeWorker{
		subscriber: subscriber,
	}
}

// Name returns the worker name.
func (w *SubscribeWorker) Name() string {
	return "subscribe-worker"
}

// Run executes the worker.
func (w *SubscribeWorker) Run(ctx context.Context) error {
	return w.subscriber.Subscribe(ctx)
}
