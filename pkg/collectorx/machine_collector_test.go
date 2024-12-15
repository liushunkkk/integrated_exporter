package collectorx

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/liushun-ing/integrated_exporter/config"
	"github.com/liushun-ing/integrated_exporter/pkg/constantx"
	"github.com/liushun-ing/integrated_exporter/pkg/metricx"
)

func TestMachineCollector_Collect(t *testing.T) {
	machineConfig := config.MachineConfig{
		Metrics:   constantx.MachineAll,
		Mounts:    []string{"/"},
		Processes: []string{"sleep"},
	}
	registry := metricx.NewIRegistry()
	mc := NewMachineCollector(machineConfig, "test_machine", registry)
	mc.Collect()

	metrics, err := registry.ExportMetrics()
	assert.Nil(t, err)

	assert.Contains(t, metrics.String(), "machine_cpu_core")
	assert.Contains(t, metrics.String(), `machine_disk_free{mountpoint="/"}`)
	assert.Contains(t, metrics.String(), "machine_memory_total")
	assert.Contains(t, metrics.String(), `machine_process_cpu_percent{processname="sleep"}`)
	assert.Contains(t, metrics.String(), "machine_process_total")
	assert.Contains(t, metrics.String(), "network_connections")
}

func TestMachineCollector_Collect_NoNetwork(t *testing.T) {
	machineConfig := config.MachineConfig{
		Metrics: constantx.MachineNoNetwork,
		Mounts:  []string{"/"},
	}
	registry := metricx.NewIRegistry()
	mc := NewMachineCollector(machineConfig, "test_machine", registry)
	mc.Collect()

	metrics, err := registry.ExportMetrics()
	assert.Nil(t, err)

	assert.Contains(t, metrics.String(), "machine_cpu_core")
	assert.Contains(t, metrics.String(), `machine_disk_free{mountpoint="/"}`)
	assert.NotContains(t, metrics.String(), `machine_disk_free{mountpoint="/dev"}`)
	assert.Contains(t, metrics.String(), "machine_memory_total")
	assert.NotContains(t, metrics.String(), `machine_process_cpu_percent{processname="sleep"}`)
	assert.Contains(t, metrics.String(), "machine_process_total")
	assert.NotContains(t, metrics.String(), "network_connections")
}
