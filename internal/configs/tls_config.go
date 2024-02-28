package configs

type tlsConfig struct {
	MTLS mtlsConfig                   `yaml:"mtls"`
	CL   map[string]certificateConfig `yaml:"cl"` // map key indicates service's server name indication (TLS SNI RFC 3546)
}

type mtlsConfig struct {
	Required  bool     `yaml:"required"`
	ClientCAs []string `yaml:"client_cas"` // list of CAs whos signatures are accepted when shown by clients
}

type certificateConfig struct {
	CertShownByFrontendToClientsMatchingSni string `yaml:"cert_shown_by_frontend_to_clients_matching_sni"`
	PrivkeyForCertShownByFrontendToClient   string `yaml:"privkey_for_cert_shown_by_frontend_to_client"`
}
