package cmd

import (
	"fmt"
	"integrated-exporter/config"

	"github.com/spf13/cobra"
)

var (
	Version string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("integrated-exporter version: ", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func SetVersion(version string) {
	Version = version
	config.C.Version = version
}
