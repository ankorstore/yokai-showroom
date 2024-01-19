package service

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
)

type Subscriber interface {
	Subscribe(ctx context.Context) error
}

type PubSubSubscriber struct {
	config *config.Config
	client *pubsub.Client
}

func NewPubSubSubscriber(config *config.Config, client *pubsub.Client) *PubSubSubscriber {
	return &PubSubSubscriber{
		config: config,
		client: client,
	}
}

func (s *PubSubSubscriber) Subscribe(ctx context.Context) error {
	subscription, err := s.subscription(ctx)
	if err != nil {
		return err
	}

	return subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		log.CtxLogger(ctx).Info().Msgf("received new message, data: %v", string(msg.Data))

		msg.Ack()
	})
}

func (s *PubSubSubscriber) subscription(ctx context.Context) (*pubsub.Subscription, error) {
	logger := log.CtxLogger(ctx)

	pubSubTopicName := s.config.GetString("modules.pubsub.topics.test")
	pubSubTopic := s.client.Topic(pubSubTopicName)

	topicExists, err := pubSubTopic.Exists(ctx)
	if err != nil {
		logger.Error().Err(err).Msgf("cannot check if topic %s exist", pubSubTopicName)

		return nil, err
	}

	if !topicExists {
		logger.Info().Msgf("topic %s does not exist, creating it", pubSubTopicName)

		pubSubTopic, err = s.client.CreateTopic(ctx, pubSubTopicName)
		if err != nil {
			logger.Error().Err(err).Msgf("cannot create topic %s", pubSubTopicName)

			return nil, err
		}
	}

	pubSubSubscriptionName := s.config.GetString("modules.pubsub.subscriptions.test")
	pubSubSubscription := s.client.Subscription(pubSubSubscriptionName)

	subscriptionExists, err := pubSubSubscription.Exists(ctx)
	if err != nil {
		logger.Error().Err(err).Msgf("cannot check if subscription %s exist", pubSubSubscriptionName)

		return nil, err
	}

	if !subscriptionExists {
		logger.Info().Msgf("subscription %s does not exist, creating it", pubSubSubscriptionName)

		pubSubSubscription, err = s.client.CreateSubscription(
			ctx,
			pubSubSubscriptionName,
			pubsub.SubscriptionConfig{
				Topic:       pubSubTopic,
				AckDeadline: 10 * time.Second,
			},
		)
		if err != nil {
			logger.Error().Err(err).Msgf("cannot create subscription %s", pubSubSubscriptionName)
			return nil, err
		}
	}

	return pubSubSubscription, nil
}
