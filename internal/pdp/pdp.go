package pdp

import (
	"log"

	"github.com/leobrada/ztsfc_proxy/internal/configs"
)

// Policy Decision Point (PDP) struct defining the main access control instance for the ZTSFC proxy
type PDP struct {
	// ControlPlane logger PDP uses for logging all its actions
	cpLogger *log.Logger
}

// NewPDP creates a new Policy Decision Point (PDP) instance using the provided configuration and logger.
// It initializes and returns the PDP instance.
// Parameters:
//   - config: A pointer to the configuration struct holding PEP settings and service configurations.
//   - dataPlaneLogger: A pointer to the logger instance for data plane logging.
//
// Returns:
//   - *PDP: A pointer to the created PDP instance.
//   - error: An error if any occurred during initialization.
func NewPDP(config *configs.Config, controlPlaneLogger *log.Logger) (*PDP, error) {
	// Create a new PEP instance with the provided logger and initialized services.
	return &PDP{
		cpLogger: controlPlaneLogger,
	}, nil
}
