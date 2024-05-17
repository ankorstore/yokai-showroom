package internal

import (
	"context"
	"fmt"
	"testing"

	"github.com/ankorstore/yokai-showroom/http-demo/db/seeds"
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/ankorstore/yokai/fxsql"
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
	fxsql.FxSQLModule,
	// dependencies registration
	Register(),
	// routing registration
	Router(),
)

// Run starts the application, with a provided [context.Context].
func Run(ctx context.Context) {
	Bootstrapper.WithContext(ctx).RunApp(
		// run SQL migrations
		fxsql.RunFxSQLMigration("up"),
	)
}

// RunTest starts the application in test mode, with an optional list of [fx.Option].
func RunTest(tb testing.TB, options ...fx.Option) {
	tb.Helper()

	// configs
	tb.Setenv("APP_CONFIG_PATH", fmt.Sprintf("%s/configs", RootDir))
	tb.Setenv("MODULES_SQL_MIGRATIONS_PATH", fmt.Sprintf("%s/db/migrations", RootDir))
	tb.Setenv("MODULES_HTTP_SERVER_TEMPLATES_ENABLED", "true")
	tb.Setenv("MODULES_HTTP_SERVER_TEMPLATES_PATH", fmt.Sprintf("%s/templates/*.html", RootDir))

	Bootstrapper.RunTestApp(
		tb,
		// run SQL migrations
		fxsql.RunFxSQLMigration("up"),
		// register seeds
		fxsql.AsSQLSeeds(seeds.NewGophersSeed),
		// apply per test options
		fx.Options(options...),
	)
}
