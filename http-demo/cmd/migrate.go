package cmd

import (
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxsql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate [command] [args...]",
	Short: "Run application database migrations",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fxcore.
			NewBootstrapper().
			WithContext(cmd.Context()).
			WithOptions(
				fx.NopLogger,
				// modules
				fxsql.FxSQLModule,
				// migrate and shutdown
				fxsql.RunFxSQLMigrationAndShutdown(args[0], args[1:]...),
			).
			RunApp()
	},
}
