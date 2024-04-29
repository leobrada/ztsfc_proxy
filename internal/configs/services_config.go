package configs

// ServicesConfig holds configurations applicable to all services managed by the Policy Enforcement Point (PEP).
// It includes a global TLS configuration to secure communications and a map of service-specific configurations.
type ServicesConfig struct {
	TLS         TLSConfig                `yaml:"tls"`          // TLS specifies the common Transport Layer Security settings applied to all services.
	ServicePool map[string]ServiceConfig `yaml:"service_pool"` // ServicePool maps service identifiers to their respective configurations.
}

// ServiceConfig defines the configuration details for a single service managed by the PEP.
// It primarily contains the URL where the service can be accessed.
type ServiceConfig struct {
	ServiceURL string `yaml:"service_url"` // ServiceURL is the endpoint URL where the service is accessible, e.g., "https://api.example.com/service".
}
