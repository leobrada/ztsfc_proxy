package service

import "fmt"

type ServicePool struct {
	Services     map[string]*Service `yaml:"services"`
	SniToService map[string]*Service
}

func (sp *ServicePool) InitServicePool() error {
	// sp.Services is initialized by LoadYAMLConfig() if it is not empty in yaml file
	if sp.Services == nil {
		return fmt.Errorf("service.InitServicePool(): service pool is empty")
	}

	sp.SniToService = make(map[string]*Service)
	for _, service := range sp.Services {
		if err := service.InitService(); err != nil {
			return fmt.Errorf("service.InitServicePool(): %v", err)
		}
		sp.SniToService[service.Sni] = service
	}

	return nil
}
