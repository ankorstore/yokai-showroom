package internal

import (
	"context"
	"fmt"
	"testing"

	"github.com/ankorstore/yokai-showroom/http-demo/internal/model"
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/ankorstore/yokai/fxorm"
	"go.uber.org/fx"
)

func init() {
	RootDir = fxcore.RootDir(1)
}

// RootDir is the application root directory.
var RootDir string

// Bootstrapper can be used to load modules, options, dependencies, routing and bootstraps the application.
var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// modules registration
	fxhttpserver.FxHttpServerModule,
	fxorm.FxOrmModule,
	// dependencies registration
	Register(),
	// routing registration
	Router(),
)

// Run starts the application, with a provided [context.Context].
func Run(ctx context.Context) {
	Bootstrapper.WithContext(ctx).RunApp(
		// run orm migrations
		fxorm.RunFxOrmAutoMigrate(&model.Gopher{}),
	)
}

// RunTest starts the application in test mode, with an optional list of [fx.Option].
func RunTest(tb testing.TB, options ...fx.Option) {
	tb.Helper()

	// configs
	tb.Setenv("APP_CONFIG_PATH", fmt.Sprintf("%s/configs", RootDir))

	// templates
	tb.Setenv("MODULES_HTTP_SERVER_TEMPLATES_ENABLED", "true")
	tb.Setenv("MODULES_HTTP_SERVER_TEMPLATES_PATH", fmt.Sprintf("%s/templates/*.html", RootDir))

	Bootstrapper.RunTestApp(
		tb,
		// run orm migrations
		fxorm.RunFxOrmAutoMigrate(&model.Gopher{}),
		// apply per test options
		fx.Options(options...),
	)
}
