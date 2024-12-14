package main

import (
	"github.com/liushun-ing/integrated_exporter/cmd"
	"github.com/liushun-ing/integrated_exporter/config"
)

const (
	APP       = "github.com/liushun-ing/integrated_exporter"
	VERSION   = "1.0.0"
	EnvPrefix = "CLI"
)

func main() {
	config.SetEnvPrefix(EnvPrefix)
	config.SetAPP(APP)
	cmd.SetVersion(VERSION)
	cmd.Execute()
}
