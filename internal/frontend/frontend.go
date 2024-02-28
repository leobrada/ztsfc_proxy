package frontend

import (
	"fmt"
	"net/http"
	"time"

	"github.com/leobrada/ztsfc_proxy/internal/configs"
	"github.com/leobrada/ztsfc_proxy/internal/logger"
	"github.com/leobrada/ztsfc_proxy/internal/pep"
	"github.com/leobrada/ztsfc_proxy/internal/security/tlsutil"
)

func NewFrontend(config *configs.Config) (*http.Server, error) {
	dpLogger, err := logger.NewDataPlaneLogger(config)
	if err != nil {
		return nil, fmt.Errorf("frontend.NewFrontend(): %v", err)
	}

	tls, err := tlsutil.NewTLS(config)
	if err != nil {
		return nil, fmt.Errorf("frontend.NewFrontend(): %v", err)
	}

	pep, err := pep.NewPEP(config, dpLogger)
	if err != nil {
		return nil, fmt.Errorf("frontend.NewFrontend(): %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", pep)

	frontend := &http.Server{
		Addr:              config.Frontend.Addr,
		Handler:           mux,
		TLSConfig:         tls,
		ReadHeaderTimeout: time.Second * 5,
		ErrorLog:          dpLogger,
	}

	return frontend, nil
}
