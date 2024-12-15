package proberx

import (
	"fmt"
	"github.com/liushun-ing/integrated_exporter/config"
	"github.com/liushun-ing/integrated_exporter/pkg/constantx"
	"github.com/shirou/gopsutil/v4/process"
	"log"
	"strings"
)

func ProbeProcess(ps config.ProcessService) error {
	processes, err := process.Processes()
	if err != nil {
		log.Printf("Failed to get process list for probe %s %s: %v", constantx.ProcessService, ps.Name, err)
		return err
	}
	for _, proc := range processes {
		cmdline, err := proc.Name()
		if err != nil {
			continue
		}
		if strings.Contains(cmdline, ps.Name) {
			return nil
		}
	}
	return fmt.Errorf("process not found for probe %s %s", constantx.ProcessService, ps.Name)
}
