package tlsutil

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	gct "github.com/leobrada/golang_convenience_tools"
	"github.com/leobrada/ztsfc_proxy/internal/configs"
	"github.com/leobrada/ztsfc_proxy/internal/logger"
)

// NewClientTLS creates a new TLS configuration for client-side connections using the provided TLS configuration.
// It initializes certificate maps used for client authentication and certificate authorities (CAs), and certificate revocation lists (CRLs) for server verification.
// Parameters:
//   - tlsConfig: A pointer to the configuration struct holding TLS settings.
//
// Returns:
//   - *tls.Config: A pointer to the created TLS configuration.
//   - error: An error if any occurred during initialization.
func NewClientTLS(tlsConfig *configs.TLSConfig) (*tls.Config, error) {
	// Initialize certificate map holding all certificates used to authenticate against a server.
	// Indexed by the target server's SNI
	cm, err := NewCertificateMap(tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("tlsutil.NewTLS(): could not load CL: %v", err)
	}

	// Initialize certificate authorities (CAs) for server verification.
	serverCAs, serverCAsListForCRLChecking, err := NewCAs(tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("tlsutil.NewServerTLS(): could not load client CA list: %v", err)
	}

	// Initialize certificate revocation list (CRL) for server certificate verification.
	serverCRL, err := NewCRL(tlsConfig, serverCAsListForCRLChecking)
	if err != nil {
		return nil, fmt.Errorf("tlsutil.NewTLS(): could not load internal CRL: %v", err)
	}

	// Create a new TLS configuration for the client.
	clientTLS := &tls.Config{
		Rand:                   nil,
		Time:                   nil,
		InsecureSkipVerify:     false,
		NextProtos:             []string{"h2"}, // enforces HTTP/2
		MinVersion:             tls.VersionTLS13,
		MaxVersion:             tls.VersionTLS13,
		SessionTicketsDisabled: true,
		Certificates:           nil,
		RootCAs:                serverCAs,
		GetClientCertificate:   makeGetClientCertificateFunction(cm),
		VerifyConnection:       makeVerifyConnection(serverCRL),
	}
	return clientTLS, nil
}

// NewServerTLS creates a new TLS configuration for server-side connections using the provided TLS configuration.
// It initializes certificate authorities (CAs) and certificate revocation lists (CRLs) for client verification.
// And certificate maps storing certificates for showing to clients, and client authentication settings.
// Parameters:
//   - tlsConfig: A pointer to the configuration struct holding TLS settings.
//
// Returns:
//   - *tls.Config: A pointer to the created TLS configuration.
//   - error: An error if any occurred during initialization.
func NewServerTLS(tlsConfig *configs.TLSConfig) (*tls.Config, error) {
	// Initialize certificate authorities (CAs) for client verification.
	// Holding the CAs that are accepted to sign client certififactes and client CRLs
	clientCAs, clientCAsListForCRLChecking, err := NewCAs(tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("tlsutil.NewServerTLS(): could not load client CA list: %v", err)
	}
	// Log the number of loaded client CAs.
	logger.SystemLogger.Debugf("tlsutil.NewServerTLS(): %d client CA(s) %s loaded", len(clientCAsListForCRLChecking), logger.Success)

	// Initialize certificate revocation list (CRL) for client certificate verification.
	clientCRL, err := NewCRL(tlsConfig, clientCAsListForCRLChecking)
	if err != nil {
		return nil, fmt.Errorf("tlsutil.NewServerTLS(): could not load external CRL: %v", err)
	}

	// Initialize certificate map storing all server certificates shown to clients for server authentication
	// Certificates are indexed by the requested service's SNI
	cm, err := NewCertificateMap(tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("tlsutil.NewServerTLS(): could not load CL: %v", err)
	}

	// Create a new TLS configuration for the server.
	serverTLS := &tls.Config{
		Rand:                   nil,
		Time:                   nil,
		InsecureSkipVerify:     false,
		NextProtos:             []string{"h2"}, // enforces HTTP/2
		MinVersion:             tls.VersionTLS13,
		MaxVersion:             tls.VersionTLS13,
		SessionTicketsDisabled: true,
		Certificates:           nil,
		ClientAuth:             setMTLS(tlsConfig),
		ClientCAs:              clientCAs,
		GetCertificate:         makeGetCertificateFunction(cm),
		VerifyConnection:       makeVerifyConnection(clientCRL),
	}
	return serverTLS, nil
}

func setMTLS(tlsConfig *configs.TLSConfig) tls.ClientAuthType {
	if !tlsConfig.ClientAuth {
		return tls.NoClientCert
	}
	return tls.RequireAndVerifyClientCert
}

// NewCAs initializes certificate authorities (CAs) from the provided TLS configuration.
// It loads the accepted CAs into a certificate pool and returns the pool along with a list of loaded certificates.
// Parameters:
//   - tlsConfig: A pointer to the configuration struct holding TLS settings.
//
// Returns:
//   - *x509.CertPool: A pointer to the created certificate pool.
//   - []*x509.Certificate: A list of loaded certificates.
//   - error: An error if any occurred during initialization.
func NewCAs(tlsConfig *configs.TLSConfig) (*x509.CertPool, []*x509.Certificate, error) {
	// Initialize a new certificate pool.
	CACertPool := x509.NewCertPool()
	// Initialize an empty list to hold loaded certificates.
	CACertList := make([]*x509.Certificate, 0)

	// Iterate through each accepted CA from the TLS configuration.
	for _, clientCA := range tlsConfig.CAs {
		// Load the certificate into the certificate pool and update the certificate list.
		var err error
		CACertList, err = gct.LoadCertificate(clientCA, CACertPool, CACertList)
		if err != nil {
			return nil, nil, fmt.Errorf("tlsutil: NewCAs(): could not load accepted CAs: '%s'", err)
		}
	}

	// Return the certificate pool and loaded certificate list.
	return CACertPool, CACertList, nil
}

// NewCRL loads and verifies a certificate revocation list (CRL) using the provided TLS configuration and CA certificate list.
// Parameters:
//   - tlsConfig: A pointer to the configuration struct holding TLS settings.
//   - cAsListForCRLChecking: A list of CA certificates used for CRL signature verification.
//
// Returns:
//   - *x509.RevocationList: A pointer to the parsed and verified CRL.
//   - error: An error if any occurred during loading or verification.
func NewCRL(tlsConfig *configs.TLSConfig, cAsListForCRLChecking []*x509.Certificate) (*x509.RevocationList, error) {
	// Read the CRL file.
	CRLBinary, err := os.ReadFile(tlsConfig.CRL)
	if err != nil {
		return nil, fmt.Errorf("tlsutil.NewCRL(): Could not load CRL '%s': '%s'", tlsConfig.CRL, err)
	}

	// Parse the CRL.
	crl, err := x509.ParseRevocationList(CRLBinary)
	if err != nil {
		return nil, fmt.Errorf("tlsutil.NewCRL(): Could not parse CRL '%s': '%s'", tlsConfig.CRL, err)
	}

	// Check if the CRL lies within the valid time period.
	now := time.Now()
	if crl.ThisUpdate.After(now) || crl.NextUpdate.Before(now) {
		return nil, fmt.Errorf("tlsutil.NewCRL(): CRL '%s' lies outside of valid time period", tlsConfig.CRL)
	}

	// Verify the signature of the CRL using the CA certificates.
	signatureVerified := false
	for _, caCert := range cAsListForCRLChecking {
		if err = crl.CheckSignatureFrom(caCert); err == nil {
			logger.SystemLogger.Debugf("tlsutil.NewCRL(): Signature for CRL '%s' %s verified by CA cert '%s'", tlsConfig.CRL, logger.Success, caCert.Subject.CommonName)
			signatureVerified = true
			break
		}
	}

	// If the signature verification fails, return an error.
	if !signatureVerified {
		return nil, fmt.Errorf("tlsutil.NewCRL(): Could not verify CRL signature: '%s'", err)
	}

	// Return the parsed and verified CRL.
	return crl, nil
}

// NewCertificateMap creates a map of TLS certificates keyed by Server Name Indication (SNI) from the provided TLS configuration.
// It loads certificates for each SNI from the specified files and returns the certificate map.
// Parameters:
//   - tlsConfig: A pointer to the configuration struct holding TLS settings.
//
// Returns:
//   - map[string]tls.Certificate: A map of TLS certificates keyed by SNI.
//   - error: An error if any occurred during loading of certificates.
func NewCertificateMap(tlsConfig *configs.TLSConfig) (map[string]tls.Certificate, error) {
	// Initialize an empty map to hold TLS certificates.
	certificateMap := make(map[string]tls.Certificate)

	// Iterate through each SNI and its corresponding certificate configuration.
	for sni, certificateConfig := range tlsConfig.Certificates {
		// Load the X.509 certificate and private key pair from the specified files.
		cert, err := tls.LoadX509KeyPair(certificateConfig.CertFile, certificateConfig.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("tlsutil.NewCertificateMap(): certificate for SNI '%s' could not be loaded: %w", sni, err)
		}
		// Store the certificate in the certificate map keyed by SNI.
		certificateMap[sni] = cert
	}

	// Return the map of TLS certificates.
	return certificateMap, nil
}

// Maker function that returns the GetCertificate() function necessary for TLS configurations.
// Has access to a list of server certificates keyed by SNI
func makeGetCertificateFunction(cm map[string]tls.Certificate) func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	// GetCertificate func(*ClientHelloInfo) (*Certificate, error)
	return func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		// use SNI map to load suitable certificate
		serverCertificate, ok := cm[hello.ServerName]
		if !ok {
			return nil, fmt.Errorf("tlsutil.GetCertificate(): could not serve a suitable certificate for %s", hello.ServerName)
		}
		return &serverCertificate, nil
	}
}

// Maker function that returns the GetCertificate() function necessary for TLS configurations.
// Has access to a list of server certificates keyed by SNI
func makeGetClientCertificateFunction(cm map[string]tls.Certificate) func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
	return func(info *tls.CertificateRequestInfo) (*tls.Certificate, error) {
		if len(info.AcceptableCAs) == 0 {
			for _, clientCert := range cm {
				return &clientCert, nil
			}
		}
		for _, caDN := range info.AcceptableCAs {
			for _, clientCert := range cm {
				certIssuerDN, err := gct.GetIssuerDNInDER(clientCert)
				if err != nil {
					return nil, fmt.Errorf("tlsutil.GetClientCertificate(): Could not retrieve Issuer DN from client certificate")
				}
				if compareDNs(caDN, certIssuerDN) {
					return &clientCert, nil
				}
			}
		}
		return nil, fmt.Errorf("tlsutil.GetClientCertificate(): Could not serve a suitable certificate for acceptable CAs")
	}
}

func compareDNs(dn1, dn2 []byte) bool {
	return bytes.Equal(dn1, dn2)
}

// makeVerifyConnection creates a function for verifying TLS connections against a certificate revocation list (CRL).
// Parameters:
//   - crl: A pointer to the parsed and verified certificate revocation list.
//
// Returns:
//   - func(tls.ConnectionState) error: A function that verifies TLS connections against the provided CRL.
func makeVerifyConnection(crl *x509.RevocationList) func(tls.ConnectionState) error {
	// Define a function for verifying TLS connections.
	return func(con tls.ConnectionState) error {
		// Check if the verified chains hold a valid client certificate.
		if len(con.VerifiedChains) == 0 || len(con.VerifiedChains[0]) == 0 {
			return fmt.Errorf("tlsutil.VerifyConnection(): error: verified chains does not hold a valid client certificate")
		}

		// Iterate through revoked certificate entries in the CRL.
		for _, revokedCertificateEntry := range crl.RevokedCertificateEntries {
			// Check if the client certificate serial number matches any revoked certificate entry.
			if con.VerifiedChains[0][0].SerialNumber.Cmp(revokedCertificateEntry.SerialNumber) == 0 {
				return fmt.Errorf("tlsutil.VerifyConnection(): client '%s' certificate is revoked", con.VerifiedChains[0][0].Subject.CommonName)
			}
		}

		// Return nil if the connection is verified successfully.
		return nil
	}
}
