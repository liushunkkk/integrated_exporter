package main

import (
	"integrated-exporter/cmd"
	"integrated-exporter/config"
)

const (
	APP       = "integrated-exporter"
	VERSION   = "1.0.0"
	EnvPrefix = "CLI"
)

func main() {
	config.SetEnvPrefix(EnvPrefix)
	config.SetAPP(APP)
	cmd.SetVersion(VERSION)
	cmd.Execute()
}
