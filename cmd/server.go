package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"

	"github.com/liushunking/integrated_exporter/config"
	"github.com/liushunking/integrated_exporter/core/server"
	"github.com/liushunking/integrated_exporter/pkg/constantx"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: `integrated-exporter server`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		interval, err := time.ParseDuration(config.C.Server.Interval)
		if err != nil {
			log.Printf("Error parsing interval '%s' for config.C.Server.Interval: %v", config.C.Server.Interval, err)
			return err
		}
		for _, service := range config.C.Server.HttpServices {
			duration, err := time.ParseDuration(service.Timeout)
			if err != nil {
				return err
			}
			if duration >= interval {
				return fmt.Errorf("%s: HttpService.Timeout should be smaller than Server.Interval", service.Name)
			}
		}
		for _, service := range config.C.Server.GrpcServices {
			duration, err := time.ParseDuration(service.Timeout)
			if err != nil {
				return err
			}
			if duration >= interval {
				return fmt.Errorf("%s: GrpcService.Timeout should be smaller than Server.Interval", service.Name)
			}
		}
		for _, service := range config.C.Server.GethServices {
			duration, err := time.ParseDuration(service.Timeout)
			if err != nil {
				return err
			}
			if duration >= interval {
				return fmt.Errorf("%s: GethService.Timeout should be smaller than Server.Interval", service.Name)
			}
		}
		for _, service := range config.C.Server.ApiServices {
			duration, err := time.ParseDuration(service.Timeout)
			if err != nil {
				return err
			}
			if duration >= interval {
				return fmt.Errorf("%s: ApiService.Timeout should be smaller than Server.Interval", service.Name)
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return server.Run(config.C.Server, nil, nil)
	},
}

func init() {
	{
		rootCmd.AddCommand(serverCmd)

		serverCmd.Flags().IntP("port", "p", 6070, "exporter server port")
		serverCmd.Flags().StringP("interval", "i", "5s", "exporter server interval for probing")
		serverCmd.Flags().StringP("route", "r", "/metrics", "exporter server metrics route")
		serverCmd.Flags().StringSlice("machineConfig.metrics", constantx.MachineAll, "the machine metrics that need to be exported.")
		serverCmd.Flags().StringSlice("machineConfig.mounts", []string{"/"}, "the mount points that need disk metrics.")
		serverCmd.Flags().StringSlice("machineConfig.processes", nil, "the processes that need detailed metrics.")
	}
}
