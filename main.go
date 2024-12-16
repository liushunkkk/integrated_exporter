package main

import (
	"github.com/liushunking/integrated_exporter/cmd"
	"github.com/liushunking/integrated_exporter/config"
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
