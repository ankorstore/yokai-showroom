package internal

import (
	"context"
	"fmt"
	"testing"

	"github.com/ankorstore/yokai-showroom/worker-demo/modules/fxpubsub"
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/ankorstore/yokai/fxworker"
	"go.uber.org/fx"
)

var RootDir string

var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// modules
	fxhttpserver.FxHttpServerModule,
	fxworker.FxWorkerModule,
	fxpubsub.FxPubSubModule,
	// routing
	ProvideRouting(),
	// services
	ProvideServices(),
)

func init() {
	RootDir = fxcore.RootDir(1)
}

func Run(ctx context.Context) {
	Bootstrapper.WithContext(ctx).RunApp()
}

func RunTest(tb testing.TB, options ...fx.Option) {
	tb.Helper()

	tb.Setenv("APP_CONFIG_PATH", fmt.Sprintf("%s/configs", RootDir))
	tb.Setenv("PUBSUB_PROJECT_ID", "worker-demo-project")

	Bootstrapper.RunTestApp(tb, fx.Options(options...))
}
