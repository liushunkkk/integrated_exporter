package server

import (
	"bytes"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"integrated-exporter/config"
	"log"
	"sync"
)

const (
	HttpService = "http"
	RpcService  = "rpc"
	GethService = "geth"
	ApiService  = "api"
)

func probeServices(serverConfig config.ServerConfig) {
	var wg sync.WaitGroup
	for _, hs := range serverConfig.HttpServices {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := probeHttp(hs)
			saveLiveGauge(HttpService, hs.Name, err)
		}()
	}
	for _, rs := range serverConfig.RpcServices {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := probeRpc(rs)
			saveLiveGauge(RpcService, rs.Name, err)
		}()
	}
	for _, gs := range serverConfig.GethServices {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := probeGeth(gs)
			saveLiveGauge(GethService, gs.Name, err)
			saveServiceMetrics(GethService, gs.Name, resp)
		}()
	}
	for _, as := range serverConfig.ApiServices {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := probeApi(as)
			saveLiveGauge(ApiService, as.Name, err)
			saveServiceMetrics(ApiService, as.Name, resp)
		}()
	}
	wg.Wait()
}

func saveLiveGauge(serviceType, serviceName string, err error) {
	liveGauge := GetOrRegisterGauge(&prometheus.GaugeOpts{
		Namespace:   serviceName,
		Name:        "live_status",
		ConstLabels: prometheus.Labels{"type": serviceType},
	})
	if liveGauge != nil {
		if err == nil {
			liveGauge.Set(1)
		} else {
			liveGauge.Set(0)
		}
	} else {
		log.Printf("Failed to set live status for %s service %s", serviceType, serviceName)
	}
}

func saveServiceMetrics(serviceType, serviceName string, metrics []byte) {
	if metrics == nil {
		log.Printf("No metrics found for %s service %s", serviceType, serviceName)
	}
	lines := bytes.Split(metrics, []byte("\n"))
	var result []byte

	for _, line := range lines {
		trimmedLine := bytes.TrimSpace(line)
		if len(trimmedLine) == 0 || bytes.HasPrefix(trimmedLine, []byte("#")) {
			result = append(result, line...)
		} else {
			index := bytes.LastIndex(line, []byte("}"))
			if index != -1 {
				label := []byte(fmt.Sprintf(`,servicename="%s"`, serviceName))
				result = append(result, line[:index]...)
				result = append(result, label...)
				result = append(result, line[index:]...)
			} else {
				spaceIndex := bytes.LastIndex(line, []byte(" "))
				if spaceIndex != -1 {
					label := []byte(fmt.Sprintf(`{servicename="%s"}`, serviceName))
					result = append(result, line[:spaceIndex]...)
					result = append(result, label...)
					result = append(result, line[spaceIndex:]...)
				} else {
					label := []byte(fmt.Sprintf("%s_", serviceName))
					result = append(result, label...)
					result = append(result, line...)
				}
			}
		}
		result = append(result, '\n')
	}
	buffer := bytes.NewBuffer(result)
	metricsHandler.AddBuffer(buffer)
}
