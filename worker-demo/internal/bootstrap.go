package internal

import (
	"context"
	"fmt"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxgcppubsub"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxworker"
	"go.uber.org/fx"
)

func init() {
	RootDir = fxcore.RootDir(1)
}

// RootDir is the application root directory.
var RootDir string

// Bootstrapper can be used to load modules, options, dependencies and bootstraps the application.
var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// modules registration
	fxworker.FxWorkerModule,
	fxgcppubsub.FxGcpPubSubModule,
	// dependencies registration
	Register(),
)

// Run starts the application, with a provided [context.Context].
func Run(ctx context.Context) {
	Bootstrapper.WithContext(ctx).RunApp()
}

// RunTest starts the application in test mode, with an optional list of [fx.Option].
func RunTest(tb testing.TB, options ...fx.Option) {
	tb.Helper()

	Bootstrapper.RunTestApp(
		tb,
		// config lookup
		fxconfig.AsConfigPath(fmt.Sprintf("%s/configs/", RootDir)),
		// test options
		fx.Options(options...),
	)
}
