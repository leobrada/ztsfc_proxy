package service

import (
	"fmt"
	"net/url"

	"github.com/leobrada/ztsfc_proxy/internal/configs"
)

type Service struct {
	ServiceUrl *url.URL
}

func NewService(serviceConf *configs.ServiceConfig) (*Service, error) {
	serviceURL, err := url.Parse(serviceConf.ServiceURL)
	if err != nil {
		return nil, fmt.Errorf("service.NewService(): %v", err)
	}
	return &Service{ServiceUrl: serviceURL}, nil
}

/*
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
*/
