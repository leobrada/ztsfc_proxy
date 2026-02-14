// Package configs defines the YAML-backed configuration structures used by the proxy.
package configs

// frontendConfig encapsulates settings for the public-facing server.
// It includes the bind address and all TLS settings for client connections.
type frontendConfig struct {
	Addr string    `yaml:"addr"` // Addr specifies the IP address and port on which the frontend server should listen. Example: "127.0.0.1:443".
	TLS  TLSConfig `yaml:"tls"`  // TLS configures the Transport Layer Security settings for the frontend server to ensure secure communication.
}
