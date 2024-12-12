package config

import (
	"encoding/json"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// C global command flags
var C = &Config{}

type Config struct {
	App     string `mapstructure:"app"`
	Syntax  string `mapstructure:"syntax"`
	Version string `mapstructure:"version"`

	Server ServerConfig `mapstructure:"server"`
}

type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	GethServices []GethService `mapstructure:"gethServices"`
	ApiServices  []ApiService  `mapstructure:"apiServices"`
	HttpServices []HttpService `mapstructure:"httpServices"`
	RpcServices  []RpcService  `mapstructure:"rpcServices"`
}

type HttpService struct {
	Name     string `mapstructure:"name"`
	Address  string `mapstructure:"address"`
	Token    string `mapstructure:"token"`
	Method   string `mapstructure:"method"`
	Body     string `mapstructure:"body"`
	Response string `mapstructure:"response"`
}

type RpcService struct {
	Name      string `mapstructure:"name"`
	Address   string `mapstructure:"address"`
	Token     string `mapstructure:"token"`
	RpcMethod string `mapstructure:"rpcMethod"`
	Body      string `mapstructure:"body"`
	Response  string `mapstructure:"response"`
}

type ApiService struct {
	Name    string `mapstructure:"name"`
	Address string `json:"address"`
	Token   string `json:"token"`
}

type GethService struct {
	Name    string `mapstructure:"name"`
	Address string `json:"address"`
	Token   string `json:"token"`
}

func SetAPP(app string) {
	C.App = app
}

func SetEnvPrefix(envPrefix string) {
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
