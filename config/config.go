package config

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"strings"
)

var Conf ConfYaml

var defaultConf = []byte(`
core:
  server_name: "vm1"
  port: 9090
  max_cpu_usage: 80
  max_mem_usage: 55
`)

type ConfYaml struct {
	Core SectionCore `yaml:"core"`
}

// SectionCore is sub section of config.
type SectionCore struct {
	ServerName  string `yaml:"server_name"`
	Port        int    `yaml:"port"`
	MaxCpuUsage int    `yaml:"max_cpu_usage"`
	MaxMemUsage int    `yaml:"max_mem_usage"`
}

// LoadConf load config from file and read in environment variables that match
func LoadConf(confPath string) (ConfYaml, error) {
	var conf ConfYaml

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()           // read in environment variables that match
	viper.SetEnvPrefix("GO_CLOUD") // will be uppercased automatically
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if confPath != "" {
		content, err := ioutil.ReadFile(confPath)

		if err != nil {
			log.Errorf("FileRepo does not exist : %s", confPath)
			return conf, err
		}

		if err := viper.ReadConfig(bytes.NewBuffer(content)); err != nil {
			return conf, err
		}
	} else {
		// Search config in home directory with name ".pkg" (without extension).
		viper.AddConfigPath("/etc/go-cloud/")
		viper.AddConfigPath("$HOME/.go-cloud")
		viper.AddConfigPath(".")
		viper.SetConfigName("config")

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		} else {
			// load default config
			if err := viper.ReadConfig(bytes.NewBuffer(defaultConf)); err != nil {
				return conf, err
			}
		}
	}

	// Core
	conf.Core.ServerName = viper.GetString("core.server_name")
	conf.Core.Port = viper.GetInt("core.port")
	conf.Core.MaxMemUsage = viper.GetInt("core.max_mem_usage")
	conf.Core.MaxCpuUsage = viper.GetInt("core.max_cpu_usage")
	return conf, nil
}

func Initialize(path string) {
	var err error
	Conf, err = LoadConf(path)
	if err != nil {
		log.Fatalf("Load yaml config file error: '%v'", err)
		return
	}
}
