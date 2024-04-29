package pep

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/leobrada/ztsfc_proxy/internal/configs"
	"github.com/leobrada/ztsfc_proxy/internal/service"
	"github.com/leobrada/ztsfc_proxy/internal/web"
)

// Policy Enforcement Point (PEP) struct defining the main HTTP handler for the frontend HTTP server
type PEP struct {
	// DataPlane logger PEP uses for logging all its actions
	dpLogger *log.Logger
	// Pointer to all services served by the PEP
	services *service.Services
}

// NewPEP creates a new Policy Enforcement Point (PEP) instance using the provided configuration and logger.
// It initializes services based on the configuration and returns the PEP instance.
// Parameters:
//   - config: A pointer to the configuration struct holding PEP settings and service configurations.
//   - dataPlaneLogger: A pointer to the logger instance for data plane logging.
//
// Returns:
//   - *PEP: A pointer to the created PEP instance.
//   - error: An error if any occurred during initialization.
func NewPEP(config *configs.Config, dataPlaneLogger *log.Logger) (*PEP, error) {
	// Initialize services based on the configuration.
	services, err := service.NewServices(&config.Services)
	if err != nil {
		return nil, fmt.Errorf("pep.NewPEP(): %v", err)
	}

	// Create a new PEP instance with the provided logger and initialized services.
	return &PEP{
		dpLogger: dataPlaneLogger,
		services: services,
	}, nil
}

func (pep *PEP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	targetSNI := r.TLS.ServerName
	targetService, ok := pep.services.ServicePool[targetSNI]
	if !ok {
		pep.dpLogger.Printf("pep.ServeHTTP(): requested service %s could not be served", targetSNI)
		web.Handle404(w)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetService.ServiceUrl)
	if proxy == nil {
		pep.dpLogger.Printf("pep.ServeHTTP(): while serving requested service %s an internal error occured", targetSNI)
		web.Handle500(w)
		return
	}

	if pep != nil {
		proxy.ErrorLog = pep.dpLogger
	}

	proxyTransport, err := GetHTTPTransportForSchemeAndTLS(targetService.ServiceUrl.Scheme, pep.services.ServicesTLS)
	if err != nil {
		pep.dpLogger.Printf("pep.ServeHTTP(): requested service %s does not implement requested scheme", targetSNI)
		web.Handle501(w)
		return
	}
	proxy.Transport = proxyTransport

	proxy.ServeHTTP(w, r)
}

// GetHTTPTransportForSchemeAndTLS returns an HTTP transport based on the provided scheme and TLS configuration.
// Parameters:
//   - scheme: The URI scheme ("http" or "https").
//   - servicesTLS: A pointer to the TLS configuration for services.
//
// Returns:
//   - *http.Transport: An HTTP transport configured based on the scheme and TLS settings.
//   - error: An error if the scheme is unsupported.
func GetHTTPTransportForSchemeAndTLS(scheme string, servicesTLS *tls.Config) (*http.Transport, error) {
	var transport *http.Transport

	switch scheme {
	case "https":
		// Create an HTTP transport with TLS configuration for HTTPS scheme.
		transport = &http.Transport{
			IdleConnTimeout:     10 * time.Second,
			MaxIdleConnsPerHost: 10000,
			TLSClientConfig:     servicesTLS,
		}
	case "http":
		// Create an HTTP transport without TLS for HTTP scheme.
		transport = &http.Transport{
			IdleConnTimeout:     10 * time.Second,
			MaxIdleConnsPerHost: 10000,
		}
	default:
		// Return an error for unsupported scheme.
		return nil, errors.New("unsupported scheme")
	}

	return transport, nil
}
