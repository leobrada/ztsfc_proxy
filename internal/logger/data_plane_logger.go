package logger

import (
	"fmt"
	"log"

	gct "github.com/leobrada/golang_convenience_tools"
	"github.com/leobrada/ztsfc_proxy/internal/configs"
)

// NewDataPlaneLogger creates a new logger instance for the data plane with the provided configuration.
// It initializes the logger output and sets up a custom logger with specific prefix and flags.
// Parameters:
//   - dpLoggerConfig: A pointer to the configuration struct holding logger settings for the data plane.
//
// Returns:
//   - *log.Logger: A pointer to the created logger instance.
//   - error: An error if any occurred during initialization.
func NewDataPlaneLogger(dpLoggerConfig *configs.LoggerConfig) (*log.Logger, error) {
	// Get the writer for logger output.
	loggerOutput, err := gct.GetWriter(dpLoggerConfig.Output)
	if err != nil {
		return nil, fmt.Errorf("logger.NewDataPlaneLogger(): %v", err)
	}

	// Create a new logger instance with the configured output, prefix, and flags.
	dpLogger := log.New(loggerOutput, "[DataPlane] - ", log.Ldate|log.Ltime)

	// Log a debug message indicating successful initialization of the Data Plane logger.
	SystemLogger.Debugf("logger.NewDataPlaneLogger(): DataPaneLogger %s initialized. Output set to '%s'", Success, dpLoggerConfig.Output)

	return dpLogger, nil
}
