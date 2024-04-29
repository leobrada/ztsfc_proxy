package configs

// frontendConfig encapsulates the configuration settings necessary for initializing and running the frontend HTTP server.
// It holds detailed settings such as the server's address and TLS configuration.
type frontendConfig struct {
	Addr string    `yaml:"addr"` // Addr specifies the IP address and port on which the frontend server should listen. Example: "127.0.0.1:443".
	TLS  TLSConfig `yaml:"tls"`  // TLS configures the Transport Layer Security settings for the frontend server to ensure secure communication.
}
