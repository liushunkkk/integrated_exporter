package server

import (
	"fmt"
	"integrated-exporter/config"
	"integrated-exporter/pkg/metricx"
	"log"
	"net/http"
	"time"
)

func Run(serverConfig config.ServerConfig) error {
	interval, err := time.ParseDuration(serverConfig.Interval)
	if err != nil {
		return fmt.Errorf("failed to parse interval: %s", err)
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	go func() {
		collect(serverConfig)
		for {
			select {
			case <-ticker.C:
				collect(serverConfig)
			}
		}
	}()

	//http.Handle("/metrics", promhttp.HandlerFor(Reg, promhttp.HandlerOpts{Registry: Reg}))
	http.Handle(serverConfig.Route, DefaultMetricsHandler)

	log.Printf("export /metrics on port :%v", serverConfig.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%v", serverConfig.Port), nil)
	if err != nil {
		return err
	}
	return nil
}

func collect(serverConfig config.ServerConfig) {
	DefaultMetricsHandler.ClearBuffer()
	probeServices(serverConfig)
	metricsText, err := metricx.ExportDefaultIRegistryMetrics()
	if err == nil {
		DefaultMetricsHandler.AddBuffer(metricsText)
	}
}
