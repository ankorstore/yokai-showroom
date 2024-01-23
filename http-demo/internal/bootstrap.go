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

var RootDir string

var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// modules
	fxhttpserver.FxHttpServerModule,
	fxorm.FxOrmModule,
	// routing
	ProvideRouting(),
	// services
	ProvideServices(),
)

func init() {
	RootDir = fxcore.RootDir(1)
}

func Run(ctx context.Context) {
	Bootstrapper.WithContext(ctx).RunApp(
		// run orm migrations
		fxorm.RunFxOrmAutoMigrate(&model.Gopher{}),
	)
}

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
