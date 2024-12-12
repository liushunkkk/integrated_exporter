package cmd

import (
	"integrated-exporter/config"
	"integrated-exporter/internal/server"

	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: `integrated-exporter server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return server.Run(config.C.Server)
	},
}

func init() {
	{
		rootCmd.AddCommand(serverCmd)

		serverCmd.Flags().StringP("port", "p", "8888", "exporter server port")
	}
}
