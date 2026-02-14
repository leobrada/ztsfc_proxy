// Package configs defines the YAML-backed configuration structures used by the proxy.
package configs

// LoggerConfig defines the output destination for the logger.
// It is shared by both data plane and control plane loggers.
type LoggerConfig struct {
	Output string `yaml:"output"` // The file descriptor the logger output is written to
}
