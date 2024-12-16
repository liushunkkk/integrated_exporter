package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/liushunking/integrated_exporter/config"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("APP name:", config.C.App)
		fmt.Println("APP version:", config.C.Version)
		fmt.Println("APP env variables prefix:", config.EnvPrefix)
		fmt.Println("APP config file syntax:", config.C.Syntax)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func SetVersion(version string) {
	config.C.Version = version
}
