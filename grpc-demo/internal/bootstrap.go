package internal

import (
	"context"
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxgrpcserver"
	"go.uber.org/fx"
)

func init() {
	RootDir = fxcore.RootDir(1)
}

// RootDir is the application root directory.
var RootDir string

// Bootstrapper can be used to load modules, options, services and bootstraps your application.
var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// modules
	fxgrpcserver.FxGrpcServerModule,
	// services
	ProvideServices(),
)

// Run starts the application, with a provided [context.Context].
func Run(ctx context.Context) {
	Bootstrapper.WithContext(ctx).RunApp()
}

// RunTest starts the application in test mode, with an optional list of [fx.Option].
func RunTest(tb testing.TB, options ...fx.Option) {
	tb.Helper()

	tb.Setenv("APP_CONFIG_PATH", fmt.Sprintf("%s/configs", RootDir))

	Bootstrapper.RunTestApp(tb, fx.Options(options...))
}
