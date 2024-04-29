package logger

import (
	"fmt"
	"log"
	"os"

	gct "github.com/leobrada/golang_convenience_tools"
	"github.com/leobrada/ztsfc_proxy/internal/configs"
)

// NewDataPlaneLogger creates a new logger instance for the Data Plane with the provided configuration.
// It initializes the logger output and sets up a custom logger with specific prefix and flags.
// Parameters:
//   - dpLoggerConfig: A pointer to the configuration struct holding logger settings for the Data Plane.
//
// Returns:
//   - *log.Logger: A pointer to the created logger instance.
//   - error: An error if any occurred during initialization.
func NewDataPlaneLogger(dpLoggerConfig *configs.LoggerConfig) (*log.Logger, error) {
	// Get the writer for logger output.
	loggerOutput, err := gct.GetWriter(dpLoggerConfig.Output)
	if err != nil {
		// If an error occurs, set the output to standard error and return the error.
		SystemLogger.SetOutput(os.Stderr)
		return nil, fmt.Errorf("logger.NewDataPlaneLogger(): %v", err)
	}

	// Create a new logger instance with the configured output, prefix, and flags.
	dpLogger := log.New(loggerOutput, "\033[31m[DataPlane]\033[0m - ", log.Ldate|log.Ltime|log.Lmicroseconds)

	// Log a debug message indicating successful initialization of the Data Plane logger.
	SystemLogger.Debugf("logger.NewSystemLogger(): DataPaneLogger %s initialized. Output set to '%s'", Success, dpLoggerConfig.Output)

	return dpLogger, nil
}
