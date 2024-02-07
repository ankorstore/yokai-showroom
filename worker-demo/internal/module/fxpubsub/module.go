package fxpubsub

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	"github.com/ankorstore/yokai/config"
	"go.uber.org/fx"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const ModuleName = "pubsub"

var FxPubSubModule = fx.Module(
	ModuleName,
	fx.Provide(
		NewFxPubSub,
	),
)

type FxPubSubParam struct {
	fx.In
	LifeCycle fx.Lifecycle
	Config    *config.Config
}

func NewFxPubSub(p FxPubSubParam) (*pubsub.Client, error) {
	ctx := context.Background()

	if p.Config.IsTestEnv() {
		return createTestClient(ctx, p.LifeCycle, p.Config)
	} else {
		return createClient(ctx, p.LifeCycle, p.Config)
	}
}

func createClient(ctx context.Context, lc fx.Lifecycle, config *config.Config) (*pubsub.Client, error) {
	// client
	client, err := pubsub.NewClient(ctx, config.GetString("modules.pubsub.project"))
	if err != nil {
		return nil, fmt.Errorf("failed to create pubsub client: %w", err)
	}

	// lifecycle
	lc.Append(fx.Hook{
		// close on stop the client
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	return client, nil
}

func createTestClient(ctx context.Context, lc fx.Lifecycle, config *config.Config) (*pubsub.Client, error) {
	// test server
	srv := pstest.NewServer()

	conn, err := grpc.Dial(srv.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	// test client
	client, err := pubsub.NewClient(ctx, config.GetString("modules.pubsub.project"), option.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create test pubsub client: %w", err)
	}

	// lifecycle
	lc.Append(fx.Hook{
		// create on start the tests topic and subscription
		OnStart: func(ctx context.Context) error {
			// test topic
			topicName := config.GetString("modules.pubsub.topic")
			topic := client.Topic(topicName)

			topicExists, err := topic.Exists(ctx)
			if err != nil {
				return fmt.Errorf("cannot check if topic %s exist: %w", topicName, err)
			}

			if !topicExists {
				topic, err = client.CreateTopic(ctx, topicName)
				if err != nil {
					return fmt.Errorf("cannot create topic %s: %w", topicName, err)
				}
			}

			// test subscription
			subscriptionName := config.GetString("modules.pubsub.subscription")
			subscription := client.Subscription(subscriptionName)

			subscriptionExists, err := subscription.Exists(ctx)
			if err != nil {
				return fmt.Errorf("cannot check if subscription %s exist: %w", subscriptionName, err)
			}

			if !subscriptionExists {
				_, err = client.CreateSubscription(
					ctx,
					subscriptionName,
					pubsub.SubscriptionConfig{
						Topic:       topic,
						AckDeadline: 10 * time.Second,
					},
				)
				if err != nil {
					return fmt.Errorf("cannot create subscription %s: %w", subscriptionName, err)
				}
			}

			return nil
		},
		// close on stop the client and the test server
		OnStop: func(ctx context.Context) error {
			err = client.Close()
			if err != nil {
				return err
			}

			return srv.Close()
		},
	})

	return client, nil
}
