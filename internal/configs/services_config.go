// Package configs defines the YAML-backed configuration structures used by the proxy.
package configs

// ServicesConfig holds configuration for outbound connections to backend services.
// It includes shared TLS settings and a per-service mapping for routing.
type ServicesConfig struct {
	TLS         TLSConfig                `yaml:"tls"`          // TLS specifies the common Transport Layer Security settings applied to all services.
	ServicePool map[string]ServiceConfig `yaml:"service_pool"` // ServicePool maps service identifiers to their respective configurations.
}

// ServiceConfig defines the configuration for a single backend service.
// The service URL is parsed at startup and used by the reverse proxy.
type ServiceConfig struct {
	ServiceURL string `yaml:"service_url"` // ServiceURL is the endpoint URL where the service is accessible, e.g., "https://api.example.com/service".
}
