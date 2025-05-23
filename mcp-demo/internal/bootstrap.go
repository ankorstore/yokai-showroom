package internal

import (
	"context"
	"fmt"
	"testing"

	"github.com/ankorstore/yokai-showroom/mcp-demo/db/seeds"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxmcpserver"
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
	fxmcpserver.FxMCPServerModule,
	fxsql.FxSQLModule,
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

	// env configs
	tb.Setenv("MODULES_SQL_MIGRATIONS_PATH", fmt.Sprintf("%s/db/migrations", RootDir))

	Bootstrapper.RunTestApp(
		tb,
		// config lookup
		fxconfig.AsConfigPath(fmt.Sprintf("%s/configs/", RootDir)),
		// run SQL migrations
		fxsql.RunFxSQLMigration("up"),
		// register seeds
		fxsql.AsSQLSeeds(seeds.NewGophersSeed),
		// apply per test options
		fx.Options(options...),
	)
}
