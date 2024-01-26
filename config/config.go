package config

import (
	"fmt"

	"github.com/leobrada/yaml_tools"
)

var Config GlobalConfig

type GlobalConfig struct {
	Frontend frontendConfig
	// SystemLogger SystemLoggerConfig `yaml:"system_logger"`
}

func (config *GlobalConfig) InitGlobalConfig(confFilePath string) error {
	err := yaml_tools.LoadYamlFile(confFilePath, &Config)
	if err != nil {
		return fmt.Errorf("config.InitGlobalConfig(): could not load yaml file: %v", err)
	}

	return nil
}

type frontendConfig struct {
	Addr string `yaml:"addr"`
}

//type SystemLoggerConfig struct {
//	Output string `yaml:"output"`
//}
