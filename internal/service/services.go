// Package service models backend services and shared TLS settings.
// It builds the service pool used by the PEP for routing.
package service

import (
	"crypto/tls"
	"fmt"

	"github.com/leobrada/ztsfc_proxy/internal/configs"
	"github.com/leobrada/ztsfc_proxy/internal/security/tlsutil"
)

type Services struct {
	ServicesTLS *tls.Config
	// Key for the ServicePool Map is the target's service SNI (extracted from http.Request.TLS.ServerName in pep.ServeHTTP).
	// Used to choose the correct Service (URL) for the ReverseProxy.
	ServicePool map[string]*Service
}

// NewServices builds a service pool and shared TLS config for outbound connections.
func NewServices(servicesConfig *configs.ServicesConfig) (*Services, error) {
	servicesTLS, err := tlsutil.NewClientTLS(&servicesConfig.TLS)
	if err != nil {
		return nil, fmt.Errorf("service.NewServices(): %v", err)
	}

	// Parse and store backend URLs keyed by SNI.
	servicePool := make(map[string]*Service)
	for sni, serviceConf := range servicesConfig.ServicePool {
		service, err := NewService(&serviceConf)
		if err != nil {
			return nil, fmt.Errorf("service.NewServices(): %v", err)
		}
		servicePool[sni] = service
	}

	return &Services{
		ServicesTLS: servicesTLS,
		ServicePool: servicePool,
	}, nil
}
