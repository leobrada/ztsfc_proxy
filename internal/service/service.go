package service

import (
	"crypto/tls"
	"fmt"
	"net/url"

	gct "github.com/leobrada/golang_convenience_tools"
)

type Service struct {
	Sni                                string `yaml:"sni"`
	TargetServiceAddr                  string `yaml:"target_service_addr"`
	TargetServiceUrl                   *url.URL
	CertShownByPepToClientsMatchingSni string `yaml:"cert_shown_by_pep_to_clients_matching_sni"`
	PrivkeyForCertShownByPepToClient   string `yaml:"privkey_for_cert_shown_by_pep_to_client"`
	X509KeyPairShownByPepToClient      tls.Certificate
	CertShownByPepToService            string `yaml:"cert_shown_by_pep_to_service"`
	PrivkeyForCertShownByPepToService  string `yaml:"privkey_for_cert_shown_by_pep_to_service"`
	X509KeyPairShownByPepToService     tls.Certificate
	CertPepAcceptsWhenShownByService   string `yaml:"cert_pep_accepts_when_shown_by_service"`
}

func (s *Service) InitService() error {
	var err error

	emptyFields := s.checkEmptyStringFields()
	if len(emptyFields) > 0 {
		return fmt.Errorf("service.InitService(): following fields are empty: %s", emptyFields)
	}
	// Load X509KeyPairs shown by pep to client
	s.X509KeyPairShownByPepToClient, err = gct.LoadX509KeyPair(s.CertShownByPepToClientsMatchingSni, s.PrivkeyForCertShownByPepToClient)
	if err != nil {
		return fmt.Errorf("service.InitService(): Could not load 'X509KeyPairShownByPepToClient' for service '%s'", s.Sni)
	}

	return nil
}

func (s *Service) checkEmptyStringFields() []string {
	emptyFields := []string{}

	if s.Sni == "" {
		emptyFields = append(emptyFields, "Sni")
	}
	if s.TargetServiceAddr == "" {
		emptyFields = append(emptyFields, "TargetServiceAddr")
	}
	if s.CertShownByPepToClientsMatchingSni == "" {
		emptyFields = append(emptyFields, "CertShownByPepToClientsMatchingSni")
	}
	if s.PrivkeyForCertShownByPepToClient == "" {
		emptyFields = append(emptyFields, "PrivkeyForCertShownByPepToClient")
	}
	if s.CertShownByPepToService == "" {
		emptyFields = append(emptyFields, "CertShownByPepToService")
	}
	if s.PrivkeyForCertShownByPepToService == "" {
		emptyFields = append(emptyFields, "PrivkeyForCertShownByPepToService")
	}
	if s.CertPepAcceptsWhenShownByService == "" {
		emptyFields = append(emptyFields, "CertPepAcceptsWhenShownByService")
	}

	return emptyFields
}
