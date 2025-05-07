package cmd

import (
	"strings"

	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxsql"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:     "migrate [command] [args...]",
	Short:   "Run application database migrations",
	Example: strings.Join(migrateExamples, "\n"),
	Args:    cobra.MinimumNArgs(1),
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

var migrateExamples = []string{
	"  migrate up                    migrate the DB to the most recent version available",
	"  migrate up-by-one             migrate the DB up by 1",
	"  migrate up-to VERSION         migrate the DB to a specific VERSION",
	"  migrate down                  roll back the version by 1",
	"  migrate down-to VERSION       roll back to a specific VERSION",
	"  migrate redo                  re-run the latest migration",
	"  migrate reset                 roll back all migrations",
	"  migrate status                dump the migration status for the current DB",
	"  migrate version               print the current version of the database",
	"  migrate create NAME [sql|go]  creates new migration file with the current timestamp",
	"  migrate fix                   apply sequential ordering to migrations",
	"  migrate validate              check migration files without running them",
}
