package server

import (
	"fmt"
	"integrated-exporter/config"
	"log"
	"net/http"
	"time"
)

func Run(serverConfig config.ServerConfig) error {
	_ = NewMetricsRegistry()
	_ = NewMetricsHandler()

	ticker := time.NewTicker(5 * time.Second)
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
	http.Handle("/metrics", metricsHandler)

	log.Printf("export /metrics on port :%v", serverConfig.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", serverConfig.Port), nil)
	if err != nil {
		return err
	}
	return nil
}

func collect(serverConfig config.ServerConfig) {
	metricsHandler.ClearBuffer()
	probeServices(serverConfig)
	metricsText, err := GetMetricsText()
	if err == nil {
		metricsHandler.AddBuffer(metricsText)
	}
}
