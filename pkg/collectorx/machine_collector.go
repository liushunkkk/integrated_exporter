package collectorx

import (
	"github.com/liushun-ing/integrated_exporter/pkg/metricx"
)

type MachineCollector struct {
	Namespace string
	Registry  *metricx.IRegistry
}

func NewMachineCollector(namespace string, registry *metricx.IRegistry) *MachineCollector {
	return &MachineCollector{
		Namespace: namespace,
		Registry:  registry,
	}
}

func (mc *MachineCollector) CollectAll() {
	// cpu

	// 内存

	// 磁盘

	// 进程
}
