// Package pep implements the Policy Enforcement Point (PEP).
// It accepts incoming requests, performs SNI-based routing, and proxies traffic to backends.
package pep

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/leobrada/ztsfc_proxy/internal/configs"
	"github.com/leobrada/ztsfc_proxy/internal/security/hashutil"
	"github.com/leobrada/ztsfc_proxy/internal/service"
	"github.com/leobrada/ztsfc_proxy/internal/web"
)

// PEP is the main HTTP handler for the frontend server.
// It holds the service pool and a data-plane logger for request tracing.
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
	// The SNI is set during the TLS handshake and used to select the backend.
	targetSNI := r.TLS.ServerName
	targetService, ok := pep.services.ServicePool[targetSNI]
	if !ok {
		pep.dpLogger.Printf("pep.ServeHTTP(): requested service %s could not be served", targetSNI)
		web.Handle404(w)
		return
	}

	// Create a reverse proxy targeting the selected backend URL.
	proxy := httputil.NewSingleHostReverseProxy(targetService.ServiceUrl)
	if proxy == nil {
		pep.dpLogger.Printf("pep.ServeHTTP(): while serving requested service %s an internal error occured", targetSNI)
		web.Handle500(w)
		return
	}
	if pep != nil && pep.dpLogger != nil {
		proxy.ErrorLog = pep.dpLogger
	}
	// Calculate a request hash to match request/response pairs in logs.
	rHash := hashutil.CalcRequestHash(r)
	proxy.ModifyResponse = pep.responseDirector(rHash)

	// Choose transport based on backend scheme and configured TLS.
	proxyTransport, err := GetHTTPTransportForSchemeAndTLS(targetService.ServiceUrl.Scheme, pep.services.ServicesTLS)
	if err != nil {
		pep.dpLogger.Printf("pep.ServeHTTP(): requested service %s does not implement requested scheme", targetSNI)
		web.Handle501(w)
		return
	}
	proxy.Transport = proxyTransport

	// Log the forward operation; actual request mutation is minimal here.
	pep.requestDirector(w, r, targetService.ServiceUrl, rHash)

	proxy.ServeHTTP(w, r)
}

// requestDirector logs the proxying action.
// The hash includes a timestamp so each request is uniquely identified.
func (pep *PEP) requestDirector(w http.ResponseWriter, r *http.Request, resource *url.URL, rHash string) {
	pep.dpLogger.Printf("http: forwarding %s request from %s to %s - [Hash:'%s']", r.Method, r.RemoteAddr, resource.String()+r.URL.String(), rHash)
}

// responseDirector applies response headers and logs the response event.
func (pep *PEP) responseDirector(rHash string) func(*http.Response) error {
	return func(resp *http.Response) error {
		setHSTSHeader(resp)
		pep.dpLogger.Printf("http: serving %s to %s - [Hash:'%s']", resp.Request.URL.String(), resp.Request.RemoteAddr, rHash)
		return nil
	}
}

func setHSTSHeader(response *http.Response) {
	response.Header.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
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
