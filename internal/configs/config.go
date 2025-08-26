package configs

import (
	"fmt"

	"github.com/leobrada/yaml_tools"
)

// Config is a central structure that encapsulates configuration settings for various components of the application.
// It aggregates multiple sub-configuration structures, each corresponding to a different component.
type Config struct {
	Frontend           frontendConfig `yaml:"frontend"`             // Configuration specific to the frontend component.
	DataPlaneLogger    LoggerConfig   `yaml:"data_plane_logger"`    // Configuration for logging within the data plane.
	ControlPlaneLogger LoggerConfig   `yaml:"control_plane_logger"` // Configuration for logging within the control plane.
	Services           ServicesConfig `yaml:"services"`             // Configuration for various services the PEP serves.
}

// NewConfig creates a new Config instance by loading configuration settings from a specified YAML file.
// It returns a pointer to the Config structure if successful, or an error if the file cannot be loaded or parsed.
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
