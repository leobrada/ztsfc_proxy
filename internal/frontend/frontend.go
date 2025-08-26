package frontend

import (
	"fmt"
	"net/http"
	"time"

	"github.com/leobrada/ztsfc_proxy/internal/configs"
	"github.com/leobrada/ztsfc_proxy/internal/logger"
	"github.com/leobrada/ztsfc_proxy/internal/pdp"
	"github.com/leobrada/ztsfc_proxy/internal/pep"
	"github.com/leobrada/ztsfc_proxy/internal/security/tlsutil"
)

// NewFrontend creates a new HTTP server instance for the frontend using the provided configuration.
// It initializes necessary components such as logger, TLS configuration, and Policy Enforcement Point (PEP).
// Parameters:
//   - config: A pointer to the configuration struct holding frontend and logging settings.
//
// Returns:
//   - *http.Server: A pointer to the created HTTP server.
//   - error: An error if any occurred during initialization.
func NewFrontend(config *configs.Config) (*http.Server, error) {
	// Initialize Data Plane logger.
	dpLogger, err := logger.NewDataPlaneLogger(&config.DataPlaneLogger)
	if err != nil {
		return nil, fmt.Errorf("frontend.NewFrontend(): %v", err)
	}

	// Initialize Control Plane logger.
	cpLogger, err := logger.NewControlPlaneLogger(&config.ControlPlaneLogger)
	if err != nil {
		return nil, fmt.Errorf("frontend.NewFrontend(): %v", err)
	}

	// Initialize TLS configuration for the server.
	tls, err := tlsutil.NewServerTLS(&config.Frontend.TLS)
	if err != nil {
		return nil, fmt.Errorf("frontend.NewFrontend(): %v", err)
	}

	// Initialize Policy Decision Point (PDP).
	pdp, err := pdp.NewPDP(config, cpLogger)
	if err != nil {
		return nil, fmt.Errorf("frontend.NewFrontend(): %v", err)
	}

	// Initialize Policy Enforcement Point (PEP).
	pep, err := pep.NewPEP(config, dpLogger)
	if err != nil {
		return nil, fmt.Errorf("frontend.NewFrontend(): %v", err)
	}

	// Create a new HTTP request multiplexer.
	mux := http.NewServeMux()
	// Register the PEP handler to serve all incoming requests.
	mux.Handle("/", pep)

	// Create the frontend HTTP server instance with configured settings.
	frontend := &http.Server{
		Addr:              config.Frontend.Addr,
		Handler:           mux,
		TLSConfig:         tls,
		ReadHeaderTimeout: time.Second * 5,
		ErrorLog:          dpLogger,
	}

	return frontend, nil
}
