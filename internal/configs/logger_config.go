package configs

// Logger configuration used for specifying the settings for data plane as well as control plane logger
type LoggerConfig struct {
	Output string `yaml:"output"` // The file descriptor the logger output is written to
}
