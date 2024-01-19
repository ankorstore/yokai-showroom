package worker

import (
	"context"

	"github.com/ankorstore/yokai-showroom/worker-demo/internal/service"
)

type SubscribeWorker struct {
	subscriber service.Subscriber
}

func NewSubscribeWorker(subscriber service.Subscriber) *SubscribeWorker {
	return &SubscribeWorker{
		subscriber: subscriber,
	}
}

func (w *SubscribeWorker) Name() string {
	return "subscribe-worker"
}

func (w *SubscribeWorker) Run(ctx context.Context) error {
	return w.subscriber.Subscribe(ctx)
}
