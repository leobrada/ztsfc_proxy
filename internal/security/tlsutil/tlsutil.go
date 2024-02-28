package tlsutil

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

	gct "github.com/leobrada/golang_convenience_tools"
	"github.com/leobrada/ztsfc_proxy/internal/configs"
)

func NewTLS(config *configs.Config) (*tls.Config, error) {

	clientCAs, err := NewClientCAs(config)
	if err != nil {
		return nil, fmt.Errorf("tlsutil.NewTLS(): could not load client CA list: %v", err)
	}

	cl, err := NewCL(config)
	if err != nil {
		return nil, fmt.Errorf("tlsutil.NewTLS(): could not load CL: %v", err)
	}

	tlsConfig := &tls.Config{
		Rand:                   nil,
		Time:                   nil,
		InsecureSkipVerify:     false,
		MinVersion:             tls.VersionTLS13,
		MaxVersion:             tls.VersionTLS13,
		SessionTicketsDisabled: true,
		Certificates:           nil,
		ClientAuth:             setMTLS(config),
		ClientCAs:              clientCAs,
		GetCertificate:         makeGetCertificateFunction(cl),
		//VerifyConnection: func(con tls.ConnectionState) error {
		//	if len(con.VerifiedChains) == 0 || len(con.VerifiedChains[0]) == 0 {
		//		return fmt.Errorf("router: VerifyConnection(): error: verified chains does not hold a valid client certificate")
		//	}
		//
		//	for _, revokedCertificateEntry := range config.Config.CRLForExt.RevokedCertificateEntries {
		//		if con.VerifiedChains[0][0].SerialNumber.Cmp(revokedCertificateEntry.SerialNumber) == 0 {
		//			return fmt.Errorf("router: VerifyConnection(): client '%s' certificate is revoked", con.VerifiedChains[0][0].Subject.CommonName)
		//		}
		//	}
		//
		//	return nil
		//},
	}
	return tlsConfig, nil
}

func setMTLS(config *configs.Config) tls.ClientAuthType {
	if !config.TLS.MTLS.Required {
		return tls.NoClientCert
	} else {
		return tls.RequireAndVerifyClientCert
	}
}

func NewClientCAs(config *configs.Config) (*x509.CertPool, error) {
	clientCACertPool := x509.NewCertPool()
	for _, clientCA := range config.TLS.MTLS.ClientCAs {
		err := gct.LoadCertificate(clientCA, clientCACertPool, nil)
		if err != nil {
			return nil, fmt.Errorf("tlsutil: NewClientCAs(): could not load accepted client CAs: '%s'", err)
		}
	}
	return clientCACertPool, nil
}

// Creates the certificate list (CL) of server certificates shown to clients during TLS handshake
// Indexed by server name indication (SNI)
func NewCL(config *configs.Config) (map[string]tls.Certificate, error) {
	cl := make(map[string]tls.Certificate, 0)
	for sni, certificateConfig := range config.TLS.CL {
		cert, err := tls.LoadX509KeyPair(certificateConfig.CertShownByFrontendToClientsMatchingSni, certificateConfig.PrivkeyForCertShownByFrontendToClient)
		if err != nil {
			return nil, fmt.Errorf("tlsutil.NewCL(): certificate for SNI '%s' could not be loaded: %w", sni, err)
		}
		cl[sni] = cert

	}
	return cl, nil
}

// Maker function that returns the GetCertificate() function necessary for TLS configurations.
// Has access to a list of server certificates indexed by SNI
func makeGetCertificateFunction(cl map[string]tls.Certificate) func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	// GetCertificate func(*ClientHelloInfo) (*Certificate, error)
	return func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		// use SNI map to load suitable certificate
		serverCertificate, ok := cl[hello.ServerName]
		if !ok {
			return nil, fmt.Errorf("tlsutil.GetCertificate(): could not serve a suitable certificate for %s", hello.ServerName)
		}
		return &serverCertificate, nil
	}
}
