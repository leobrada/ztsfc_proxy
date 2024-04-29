package configs

type TLSConfig struct {
	// For server side Certificates stores certificates shown by the server to the client
	// For client side Certificates stores certificates shown by client to the server
	// map key indicates service's server name indication (TLS SNI RFC 3546)
	Certificates map[string]certificateConfig `yaml:"certificates"`
	ClientAuth   bool                         `yaml:"client_auth"`
	// list of CAs whos signatures are accepted when shown by clients
	CAs []string `yaml:"cas"`
	// certificate revocation list checked for client certificates provided by a client
	CRL string `yaml:"crl"`
}

type certificateConfig struct {
	// Specifies path to certificate
	CertFile string `yaml:"cert_file"`
	// Specifies path to private key belonging to the specified certificate
	KeyFile string `yaml:"key_file"`
}
