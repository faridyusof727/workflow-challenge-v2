package cmd

import (
	"log"
	"workflow-code-test/api/api"
	"workflow-code-test/api/pkg/di"

	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "api is a command line interface for workflow-challenge-v2",
	Run: func(cmd *cobra.Command, args []string) {
		di := di.NewService()

		container := di.Container(cmd.Context())
		if container == nil {
			log.Fatal("Failed to create container")
		}

		defer di.Shutdown(cmd.Context())

		api.NewServer(container).Start()
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}
