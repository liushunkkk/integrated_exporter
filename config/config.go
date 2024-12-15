package config

import (
	"encoding/json"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var EnvPrefix string

// C global command flags
var C = &Config{}

type Config struct {
	App     string `mapstructure:"app"`
	Syntax  string `mapstructure:"syntax"`
	Version string `mapstructure:"version"`

	Server ServerConfig `mapstructure:"server"`
}

type ServerConfig struct {
	Port            int              `mapstructure:"port"`
	Interval        string           `mapstructure:"interval"`
	Route           string           `mapstructure:"route"`
	MachineConfig   MachineConfig    `mapstructure:"machineConfig"`
	GethServices    []GethService    `mapstructure:"gethServices"`
	ApiServices     []ApiService     `mapstructure:"apiServices"`
	HttpServices    []HttpService    `mapstructure:"httpServices"`
	GrpcServices    []GrpcService    `mapstructure:"grpcServices"`
	ProcessServices []ProcessService `mapstructure:"processServices"`
}

type MachineConfig struct {
	Metrics   []string `mapstructure:"metrics"`
	Mounts    []string `mapstructure:"mounts"`
	Processes []string `mapstructure:"processes"`
}

type HttpService struct {
	Name     string `mapstructure:"name"`
	Address  string `mapstructure:"address"`
	Token    string `mapstructure:"token"`
	Method   string `mapstructure:"method"`
	Body     string `mapstructure:"body"`
	Response string `mapstructure:"response"`
	Timeout  string `mapstructure:"timeout"`
}

type GrpcService struct {
	Name      string `mapstructure:"name"`
	Address   string `mapstructure:"address"`
	Token     string `mapstructure:"token"`
	RpcMethod string `mapstructure:"rpcMethod"`
	Body      string `mapstructure:"body"`
	Response  string `mapstructure:"response"`
	Timeout   string `mapstructure:"timeout"`
}

type ApiService struct {
	Name    string `mapstructure:"name"`
	Address string `json:"address"`
	Token   string `json:"token"`
	Timeout string `mapstructure:"timeout"`
}

type GethService struct {
	Name    string `mapstructure:"name"`
	Address string `json:"address"`
	Token   string `json:"token"`
	Timeout string `mapstructure:"timeout"`
}

type ProcessService struct {
	Name   string `mapstructure:"name"`
	Target string `mapstructure:"target"`
}

func SetAPP(app string) {
	C.App = app
}

func SetEnvPrefix(envPrefix string) {
	EnvPrefix = envPrefix
	viper.SetEnvPrefix(envPrefix)
}

func SetConfig() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	err := viper.Unmarshal(C)
	cobra.CheckErr(err)
}

func (c *Config) String() string {
	marshal, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return string(marshal)
}
