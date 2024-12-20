package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/liushunkkk/integrated_exporter/config"
	"github.com/liushunkkk/integrated_exporter/pkg/metricx"
)

// Run start the integrated exporter. If registry or handler is nil, it uses the default one.
func Run(config config.ServerConfig, registry *metricx.IRegistry, handler *MetricsHandler) error {
	if registry == nil {
		registry = metricx.DefaultIRegistry
	}
	if handler == nil {
		handler = DefaultMetricsHandler
	}
	interval, err := time.ParseDuration(config.Interval)
	if err != nil {
		return fmt.Errorf("failed to parse interval: %s", err)
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	go func() {
		Collect(config, registry, handler)
		for range ticker.C {
			Collect(config, registry, handler)
		}
	}()

	http.Handle(config.Route, handler)

	log.Printf("export %s on port :%v", config.Route, config.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%v", config.Port), nil)
	if err != nil {
		return err
	}
	return nil
}

func Collect(serverConfig config.ServerConfig, registry *metricx.IRegistry, handler *MetricsHandler) {
	handler.ClearBuffer()
	ProbeServices(serverConfig, registry, handler)
	metricsText, err := registry.ExportMetrics()
	if err == nil {
		handler.AddBuffer(metricsText)
	}
}
