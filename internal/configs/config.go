package configs

import (
	"fmt"

	"github.com/leobrada/yaml_tools"
)

type Config struct {
	//Frontend        frontendConfig `yaml:"frontend"`
	//DataPlaneLogger loggerConfig   `yaml:"data_plane_logger"`
}

func NewConfig(confFilePath string) (*Config, error) {
	config := new(Config)

	err := yaml_tools.LoadYamlFileGeneric(confFilePath, config)
	if err != nil {
		return nil, fmt.Errorf("configs.InitConfig(): could not load yaml file: %v", err)
	}

	return config, nil
}
