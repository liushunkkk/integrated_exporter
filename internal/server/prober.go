package server

import (
	"bytes"
	"fmt"
	"log"
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/liushun-ing/integrated_exporter/config"
	"github.com/liushun-ing/integrated_exporter/pkg/constantx"
	"github.com/liushun-ing/integrated_exporter/pkg/metricx"
	"github.com/liushun-ing/integrated_exporter/pkg/proberx"
)

func probeServices(serverConfig config.ServerConfig) {
	var wg sync.WaitGroup
	for _, hs := range serverConfig.HttpServices {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := proberx.ProbeHttp(hs)
			saveLiveGauge(constantx.HttpService, hs.Name, err)
		}()
	}
	for _, rs := range serverConfig.RpcServices {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := proberx.ProbeRpc(rs)
			saveLiveGauge(constantx.RpcService, rs.Name, err)
		}()
	}
	for _, gs := range serverConfig.GethServices {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := proberx.ProbeGeth(gs)
			saveLiveGauge(constantx.GethService, gs.Name, err)
			saveServiceMetrics(constantx.GethService, gs.Name, resp)
		}()
	}
	for _, as := range serverConfig.ApiServices {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := proberx.ProbeApi(as)
			saveLiveGauge(constantx.ApiService, as.Name, err)
			saveServiceMetrics(constantx.ApiService, as.Name, resp)
		}()
	}
	wg.Wait()
}

func saveLiveGauge(serviceType, serviceName string, err error) {
	liveGauge := metricx.GetOrRegisterIGauge(&metricx.IOpts{
		Namespace: serviceName,
		Name:      "live_status",
		Labels:    prometheus.Labels{"type": serviceType, "servicename": serviceName},
	}, nil)
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
	DefaultMetricsHandler.AddBuffer(buffer)
}
