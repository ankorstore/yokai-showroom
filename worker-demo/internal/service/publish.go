package service

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/prometheus/client_golang/prometheus"
)

// PublishCounter is a metrics counter for published messages.
var PublishCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "worker_demo_app_messages_published_total",
	Help: "Total number of published messages",
})

// Publisher is the interface for publishers.
type Publisher interface {
	Publish(ctx context.Context, message string) error
}

// DefaultPublisher is the default [Publisher] implementation.
type DefaultPublisher struct {
	config *config.Config
	client *pubsub.Client
}

// NewDefaultPublisher creates a new [DefaultPublisher].
func NewDefaultPublisher(config *config.Config, client *pubsub.Client) *DefaultPublisher {
	return &DefaultPublisher{
		config: config,
		client: client,
	}
}

// Publish handles a message publication.
func (p *DefaultPublisher) Publish(ctx context.Context, message string) error {
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

	PublishCounter.Inc()

	return nil
}

func (p *DefaultPublisher) topic(ctx context.Context) (*pubsub.Topic, error) {
	topicName := p.config.GetString("modules.pubsub.topics.test")
	topic := p.client.Topic(topicName)

	exists, err := topic.Exists(ctx)
	if err != nil {
		log.CtxLogger(ctx).Error().Err(err).Msg("cannot check if topic exist")

		return nil, err
	}

	if !exists {
		log.CtxLogger(ctx).Info().Msgf("topic %s does not exist, creating it", topicName)

		topic, err = p.client.CreateTopic(ctx, topicName)
		if err != nil {
			log.CtxLogger(ctx).Error().Err(err).Msgf("cannot create topic %s", topicName)

			return nil, err
		}
	}

	return topic, nil
}
