package internal

import (
	"context"
	"fmt"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxgcppubsub"
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxworker"
	"go.uber.org/fx"
)

var RootDir string

var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// module
	fxworker.FxWorkerModule,
	fxgcppubsub.FxGcpPubSubModule,
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

	Bootstrapper.RunTestApp(tb, fx.Options(options...))
}
