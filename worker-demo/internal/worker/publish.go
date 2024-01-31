package worker

import (
	"context"
	"time"

	"github.com/ankorstore/yokai-showroom/worker-demo/internal/service"
	"github.com/ankorstore/yokai/log"
	"go.uber.org/fx"
)

type PublishWorker struct {
	publisher      service.Publisher
	messageChannel chan string
}

type PublishWorkerParam struct {
	fx.In
	Publisher      service.Publisher
	MessageChannel chan string `name:"pub-sub-message-channel"`
}

func NewPublishWorker(p PublishWorkerParam) *PublishWorker {
	return &PublishWorker{
		publisher:      p.Publisher,
		messageChannel: p.MessageChannel,
	}
}

func (w *PublishWorker) Name() string {
	return "publish-worker"
}

func (w *PublishWorker) Run(ctx context.Context) error {
	logger := log.CtxLogger(ctx)

	for message := range w.messageChannel {
		logger.Info().Msgf("waiting before publishing message: %s", message)

		// some sleep to have time to see processing delay in logs
		time.Sleep(3 * time.Second)

		logger.Info().Msgf("publishing message: %s", message)

		err := w.publisher.Publish(ctx, message)
		if err != nil {
			logger.Error().Err(err).Msgf("message publication error")

			return err
		}

		logger.Info().Msg("message publication success")
	}

	return nil
}
