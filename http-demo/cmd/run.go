package cmd

import (
	"github.com/ankorstore/yokai-showroom/http-demo/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run application",
	Run: func(cmd *cobra.Command, args []string) {
		internal.Run(cmd.Context())
	},
}
