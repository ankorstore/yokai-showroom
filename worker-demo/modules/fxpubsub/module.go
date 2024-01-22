package fxpubsub

import (
	"context"
	"fmt"

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
	var client *pubsub.Client
	var err error

	if p.Config.IsTestEnv() {
		client, err = createTestClient(p)
	} else {
		client, err = createClient(p)
	}

	return client, err
}

func createClient(p FxPubSubParam) (*pubsub.Client, error) {
	client, err := pubsub.NewClient(context.Background(), p.Config.GetString("modules.pubsub.project.id"))
	if err != nil {
		return nil, fmt.Errorf("failed to create pubsub client: %w", err)
	}

	return client, nil
}

func createTestClient(p FxPubSubParam) (*pubsub.Client, error) {
	srv := pstest.NewServer()

	conn, err := grpc.Dial(srv.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client, err := pubsub.NewClient(
		context.Background(),
		p.Config.GetString("modules.pubsub.project.id"),
		option.WithGRPCConn(conn),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create fake pubsub client: %w", err)
	}

	/*	p.LifeCycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
						err = srv.Close()
						err = conn.Close()
						err = client.Close()

			return nil
		},
	})*/

	return client, nil
}
