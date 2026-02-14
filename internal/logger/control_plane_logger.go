package logger

import (
	"fmt"
	"log"

	gct "github.com/leobrada/golang_convenience_tools"
	"github.com/leobrada/ztsfc_proxy/internal/configs"
)

// NewControlPlaneLogger creates a logger intended for control-plane operations.
// This is currently unused but reserved for future PDP integration.
// Parameters:
//   - cpLoggerConfig: A pointer to the configuration struct holding logger settings for the control.
//
// Returns:
//   - *log.Logger: A pointer to the created logger instance.
//   - error: An error if any occurred during initialization.
func NewControlPlaneLogger(cpLoggerConfig *configs.LoggerConfig) (*log.Logger, error) {
	// Get the writer for logger output.
	loggerOutput, err := gct.GetWriter(cpLoggerConfig.Output)
	if err != nil {
		return nil, fmt.Errorf("logger.NewControlPlaneLogger(): %v", err)
	}

	// Create a new logger instance with the configured output, prefix, and flags.
	cpLogger := log.New(loggerOutput, "[ControlPlane] - ", log.Ldate|log.Ltime)

	// Log a debug message indicating successful initialization of the Data Plane logger.
	SystemLogger.Debugf("logger.NewControlPlaneLogger(): ControlPaneLogger %s initialized. Output set to '%s'", Success, cpLoggerConfig.Output)

	return cpLogger, nil
}
