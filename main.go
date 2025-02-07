package main

import (
	"github.com/liushunkkk/integrated_exporter/cmd"
	"github.com/liushunkkk/integrated_exporter/config"
)

const (
	APP       = "integrated_exporter"
	VERSION   = "1.0.0"
	EnvPrefix = "INTEGRATEDEXPORTER"
)

func main() {
	config.SetEnvPrefix(EnvPrefix)
	config.SetAPP(APP)
	cmd.SetVersion(VERSION)
	cmd.Execute()
}
