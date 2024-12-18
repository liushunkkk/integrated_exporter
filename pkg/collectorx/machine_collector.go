package collectorx

import (
	"strconv"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"

	"github.com/liushunkkk/integrated_exporter/config"
	"github.com/liushunkkk/integrated_exporter/pkg/constantx"
	"github.com/liushunkkk/integrated_exporter/pkg/metricx"
	"github.com/liushunkkk/integrated_exporter/pkg/stringx"
)

type MachineCollector struct {
	Cfg       config.MachineConfig
	Namespace string
	Registry  *metricx.IRegistry
}

func NewMachineCollector(cfg config.MachineConfig, namespace string, registry *metricx.IRegistry) *MachineCollector {
	return &MachineCollector{
		Cfg:       cfg,
		Namespace: namespace,
		Registry:  registry,
	}
}

func (mc *MachineCollector) Collect() {
	for _, m := range mc.Cfg.Metrics {
		switch m {
		case constantx.MachineCpu:
			mc.CollectCpuMetrics()
		case constantx.MachineMemory:
			mc.CollectMemoryMetrics()
		case constantx.MachineDisk:
			mc.CollectDiskMetrics()
		case constantx.MachineProcess:
			mc.CollectProcessMetrics()
		case constantx.MachineNetwork:
			mc.CollectNetworkMetrics()
		}
	}
}

// CollectCpuMetrics collect and register cpu metrics.
func (mc *MachineCollector) CollectCpuMetrics() {
	// 获取 CPU 核心数, false 表示物理 CPU 核心数
	count, err := cpu.Counts(false)
	if err == nil {
		mc.AddGauge("cpu_core", nil, float64(count), "physical cores")
	}

	// 获取每个 CPU 的信息
	info, err := cpu.Info()
	if err == nil {
		for _, cpuInfo := range info {
			mc.AddGauge("cpu_ghz",
				map[string]string{"modelname": stringx.FilterAlphanumeric(cpuInfo.ModelName)},
				float64(cpuInfo.Mhz)/1000, "")
			break // 只取第一个即可
		}
	}

	// 获取 CPU 使用情况
	percent, err := cpu.Percent(0, true) // 0 表示没有等待时间，true 表示每个逻辑核心的使用情况
	if err == nil {
		for i, p := range percent {
			mc.AddGauge("cpu_use_percent", map[string]string{"cpuid": strconv.Itoa(i)}, p, "")
		}
	}

	percent, err = cpu.Percent(0, false)
	if err == nil {
		for _, p := range percent {
			mc.AddGauge("cpu_use_percent_all", nil, p, "avg cpu use percent")
			break // 只取第一个即可
		}
	}
}

// CollectMemoryMetrics collect and register memory and swap memory metrics.
func (mc *MachineCollector) CollectMemoryMetrics() {
	// 获取内存的总量、已用量、空闲量等信息
	vmStat, err := mem.VirtualMemory()
	if err == nil {
		mc.AddGauge("memory_total", nil, float64(vmStat.Total), "")
		mc.AddGauge("memory_used", nil, float64(vmStat.Used), "")
		mc.AddGauge("memory_free", nil, float64(vmStat.Free), "")
		mc.AddGauge("memory_available", nil, float64(vmStat.Available), "")
		mc.AddGauge("memory_used_percent", nil, vmStat.UsedPercent, "")
	}

	// 获取交换内存（Swap）信息
	swapStat, err := mem.SwapMemory()
	if err == nil {
		mc.AddGauge("swap_memory_total", nil, float64(swapStat.Total), "")
		mc.AddGauge("wap_memory_used", nil, float64(swapStat.Used), "")
		mc.AddGauge("swap_memory_free", nil, float64(swapStat.Free), "")
		mc.AddGauge("swap_memory_used_percent", nil, swapStat.UsedPercent, "")
	}
}

// CollectDiskMetrics collect and register disk metrics.
func (mc *MachineCollector) CollectDiskMetrics() {
	for _, mount := range mc.Cfg.Mounts {
		diskStat, err := disk.Usage(mount) // 挂载点
		if err == nil {
			mc.AddGauge("disk_total", map[string]string{"mountpoint": mount}, float64(diskStat.Total), "")
			mc.AddGauge("disk_used", map[string]string{"mountpoint": mount}, float64(diskStat.Used), "")
			mc.AddGauge("disk_free", map[string]string{"mountpoint": mount}, float64(diskStat.Free), "")
			mc.AddGauge("disk_used_percent", map[string]string{"mountpoint": mount}, diskStat.UsedPercent, "")
		}
	}
	counters, err := disk.IOCounters()
	if err == nil {
		var totalReadCount, totalWriteCount, totalReadBytes, totalWriteBytes uint64
		for _, counter := range counters {
			totalReadCount += counter.ReadCount
			totalReadBytes += counter.ReadBytes
			totalWriteCount += counter.WriteCount
			totalWriteBytes += counter.WriteBytes
		}
		mc.AddGauge("disk_read_count", nil, float64(totalReadCount), "")
		mc.AddGauge("disk_write_count", nil, float64(totalWriteCount), "")
		mc.AddGauge("disk_read", nil, float64(totalReadBytes), "read bytes")
		mc.AddGauge("disk_write", nil, float64(totalWriteBytes), "write bytes")
	}
}

// CollectProcessMetrics collect and register process metrics.
func (mc *MachineCollector) CollectProcessMetrics() {
	processes, err := process.Processes()
	if err != nil {
		return
	}

	total := 0
	threads := 0
	for _, p := range processes {
		total++
		numThreads, err := p.NumThreads()
		if err == nil {
			threads += int(numThreads)
		}
		name, err := p.Name()
		if err != nil {
			continue
		}
		if stringx.FuzzyContains(mc.Cfg.Processes, name) {
			// 获取进程的内存使用情况
			memoryInfo, err := p.MemoryInfo()
			if err == nil {
				mc.AddGauge("process_memory_usage", map[string]string{"processname": stringx.FilterAlphanumeric(name)}, float64(memoryInfo.RSS), "")
			}
			// 获取进程的 CPU 使用情况
			cpuPercent, err := p.CPUPercent()
			if err == nil {
				mc.AddGauge("process_cpu_percent", map[string]string{"processname": stringx.FilterAlphanumeric(name)}, cpuPercent, "")
			}
			// IO 读写
			counters, err := p.IOCounters()
			if err == nil {
				mc.AddGauge("process_read_count", map[string]string{"processname": stringx.FilterAlphanumeric(name)}, float64(counters.ReadCount), "")
				mc.AddGauge("process_write_count", map[string]string{"processname": stringx.FilterAlphanumeric(name)}, float64(counters.WriteCount), "")
				mc.AddGauge("process_read", map[string]string{"processname": stringx.FilterAlphanumeric(name)}, float64(counters.ReadBytes), "read bytes")
				mc.AddGauge("process_write", map[string]string{"processname": stringx.FilterAlphanumeric(name)}, float64(counters.WriteBytes), "write bytes")
			}
		}
	}
	mc.AddGauge("process_total", nil, float64(total), "")
	mc.AddGauge("process_thread", nil, float64(threads), "")
}

// CollectNetworkMetrics collect and register network metrics.
func (mc *MachineCollector) CollectNetworkMetrics() {
	tcps, err := net.Connections("tcp")
	if err == nil {
		mc.AddGauge("network_connections", map[string]string{"kind": "tcp"}, float64(len(tcps)), "")
	}
	udps, err := net.Connections("udp")
	if err == nil {
		mc.AddGauge("network_connections", map[string]string{"kind": "udp"}, float64(len(udps)), "")
	}
	all, err := net.Connections("all")
	if err == nil {
		mc.AddGauge("network_connections", map[string]string{"kind": "all"}, float64(len(all)), "")
	}

	// 获取网络 I/O 统计信息
	ioStats, err := net.IOCounters(false) // 参数为 true 时，返回每个接口的统计；false 则返回总计
	if err == nil {
		for _, io := range ioStats {
			mc.AddGauge("network_sent", nil, float64(io.BytesSent), "send bytes")
			mc.AddGauge("network_recv", nil, float64(io.BytesRecv), "recv bytes")
			mc.AddGauge("network_error", nil, float64(io.Errin+io.Errout), "")
			break
		}
	}
}

func (mc *MachineCollector) AddGauge(name string, labels map[string]string, value float64, help string) {
	g := metricx.GetOrRegisterIGauge(&metricx.IOpts{
		Namespace: mc.Namespace,
		Name:      name,
		Help:      help,
		Labels:    labels,
	}, mc.Registry)
	if g != nil {
		g.Set(value)
	}
}
