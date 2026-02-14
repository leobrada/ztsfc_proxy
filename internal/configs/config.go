// Package configs defines the YAML-backed configuration structures used by the proxy.
// It provides a single entrypoint for loading and validating configuration data.
package configs

import (
	"fmt"

	"github.com/leobrada/yaml_tools"
)

// Config is the root configuration structure for the proxy.
// It aggregates sub-configs for the frontend, loggers, and backend services.
type Config struct {
	Frontend           frontendConfig `yaml:"frontend"`             // Configuration specific to the frontend component.
	DataPlaneLogger    LoggerConfig   `yaml:"data_plane_logger"`    // Configuration for logging within the data plane.
	ControlPlaneLogger LoggerConfig   `yaml:"control_plane_logger"` // Configuration for logging within the control plane.
	Services           ServicesConfig `yaml:"services"`             // Configuration for various services the PEP serves.
}

// NewConfig loads configuration settings from a YAML file into a Config struct.
// It returns an error when the file cannot be read or parsed.
//
// Parameters:
//   - confFilePath: The path to the YAML configuration file.
//
// Returns:
//   - *Config: A pointer to the successfully created Config structure.
//   - error: An error message detailing any issues encountered during file loading or parsing, or nil if no issues occurred.
func NewConfig(confFilePath string) (*Config, error) {
	config := new(Config)

	// LoadYamlFileGeneric attempts to load and parse the YAML file into the config structure.
	err := yaml_tools.LoadYamlFileGeneric(confFilePath, config)
	if err != nil {
		return nil, fmt.Errorf("configs.InitConfig(): could not load yaml file: %v", err)
	}

	return config, nil
}
