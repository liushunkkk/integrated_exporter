package server

import (
	"bytes"
	"fmt"
	"log"
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/liushun-ing/integrated_exporter/config"
	"github.com/liushun-ing/integrated_exporter/pkg/collectorx"
	"github.com/liushun-ing/integrated_exporter/pkg/constantx"
	"github.com/liushun-ing/integrated_exporter/pkg/metricx"
	"github.com/liushun-ing/integrated_exporter/pkg/proberx"
)

func ProbeServices(config config.ServerConfig, registry *metricx.IRegistry, handler *MetricsHandler) {
	var wg sync.WaitGroup
	for _, hs := range config.HttpServices {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := proberx.ProbeHttp(hs)
			SaveLiveGauge(constantx.HttpService, hs.Name, err, registry)
		}()
	}
	for _, rs := range config.GrpcServices {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := proberx.ProbeGrpc(rs)
			SaveLiveGauge(constantx.GrpcService, rs.Name, err, registry)
		}()
	}
	for _, ps := range config.ProcessServices {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := proberx.ProbeProcess(ps)
			SaveLiveGauge(constantx.ProcessService, ps.Name, err, registry)
		}()
	}
	for _, gs := range config.GethServices {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := proberx.ProbeGeth(gs)
			SaveLiveGauge(constantx.GethService, gs.Name, err, registry)
			SaveServiceMetrics(constantx.GethService, gs.Name, resp, handler)
		}()
	}
	for _, as := range config.ApiServices {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := proberx.ProbeApi(as)
			SaveLiveGauge(constantx.ApiService, as.Name, err, registry)
			SaveServiceMetrics(constantx.ApiService, as.Name, resp, handler)
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		SaveMachineMetrics(config.MachineConfig, constantx.MachineService, registry)
	}()
	wg.Wait()
}

func SaveLiveGauge(serviceType, serviceName string, err error, registry *metricx.IRegistry) {
	liveGauge := metricx.GetOrRegisterIGauge(&metricx.IOpts{
		Namespace: serviceName,
		Name:      "live_status",
		Labels:    prometheus.Labels{"type": serviceType, "servicename": serviceName},
	}, registry)
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

// SaveServiceMetrics add an <servicename> tag to the metrics and save them to the metrics center.
func SaveServiceMetrics(serviceType, serviceName string, metrics []byte, handler *MetricsHandler) {
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
	handler.AddBuffer(buffer)
}

func SaveMachineMetrics(cfg config.MachineConfig, serviceType string, registry *metricx.IRegistry) {
	machineCollector := collectorx.NewMachineCollector(cfg, serviceType, registry)
	machineCollector.Collect()
}
