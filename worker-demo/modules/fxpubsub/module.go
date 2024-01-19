package fxpubsub

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/ankorstore/yokai/config"
	"go.uber.org/fx"
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
	// client
	client, err := pubsub.NewClient(context.Background(), p.Config.GetString("modules.pubsub.project.id"))
	if err != nil {
		return nil, fmt.Errorf("failed to create pubsub client: %w", err)
	}

	// lifecycle
	p.LifeCycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	return client, nil
}
