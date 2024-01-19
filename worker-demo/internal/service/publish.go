package service

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
)

type Publisher interface {
	Publish(ctx context.Context, message string) error
}

type PubSubPublisher struct {
	config *config.Config
	client *pubsub.Client
}

func NewPubSubPublisher(config *config.Config, client *pubsub.Client) *PubSubPublisher {
	return &PubSubPublisher{
		config: config,
		client: client,
	}
}

func (p *PubSubPublisher) Publish(ctx context.Context, message string) error {
	topic, err := p.topic(ctx)
	if err != nil {
		return err
	}

	pubSubMessage := &pubsub.Message{
		Data: []byte(message),
	}

	if _, err = topic.Publish(ctx, pubSubMessage).Get(ctx); err != nil {
		log.CtxLogger(ctx).Error().Err(err).Msg("cannot publish message")

		return err
	}

	return nil
}

func (h *PubSubPublisher) topic(ctx context.Context) (*pubsub.Topic, error) {
	topicName := h.config.GetString("modules.pubsub.topics.test")
	topic := h.client.Topic(topicName)

	exists, err := topic.Exists(ctx)
	if err != nil {
		log.CtxLogger(ctx).Error().Err(err).Msg("cannot check if topic exist")

		return nil, err
	}

	if !exists {
		log.CtxLogger(ctx).Info().Msgf("topic %s does not exist, creating it", topicName)

		topic, err = h.client.CreateTopic(ctx, topicName)
		if err != nil {
			log.CtxLogger(ctx).Error().Err(err).Msg("cannot create topic")

			return nil, err
		}
	}

	return topic, nil
}
