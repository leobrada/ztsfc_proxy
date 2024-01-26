package router

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/leobrada/ztsfc_proxy/config"
	"github.com/leobrada/ztsfc_proxy/internal/logger"
)

type Router struct {
	tlsConfig *tls.Config
	frontend  *http.Server
}

func NewRouter() *Router {
	router := new(Router)
	//rh := new(ResourceHandler)

	router.tlsConfig = &tls.Config{
		Rand:                   nil,
		Time:                   nil,
		InsecureSkipVerify:     false,
		MinVersion:             tls.VersionTLS13,
		MaxVersion:             tls.VersionTLS13,
		SessionTicketsDisabled: true,
		Certificates:           nil,
		//ClientAuth:             tls.RequireAndVerifyClientCert,
		//ClientCAs:              config.Config.CAcertPoolPepAcceptsFromExt,
		//GetCertificate: func(cli *tls.ClientHelloInfo) (*tls.Certificate, error) {
		//	// use SNI map to load suitable certificate
		//	service, ok := config.Config.ServiceSniMap[cli.ServerName]
		//	if !ok {
		//		return nil, fmt.Errorf("router: GetCertificate(): could not serve a suitable certificate for %s", cli.ServerName)
		//	}
		//	return &service.X509KeyPairShownByPepToClient, nil
		//},
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

	// Frontend Handlers
	mux := http.NewServeMux()
	mux.Handle("/", router)
	//mux.Handle("/9af1ecf7/", rh)

	// Setting Up the Frontend Server
	router.frontend = &http.Server{
		//Addr:         config.Config.Pep.ListenAddr,
		Addr:         config.Config.Frontend.Addr,
		TLSConfig:    router.tlsConfig,
		ReadTimeout:  time.Hour * 1,
		WriteTimeout: time.Hour * 1,
		Handler:      mux,
		ErrorLog:     log.New(logger.SystemLogger.Out, "", 0),
	}

	return router
}

func (router *Router) ListenAndServe() error {
	return router.frontend.ListenAndServe()
}

//func (router *Router) ListenAndServeTLS() error {
//	return router.frontend.ListenAndServeTLS("", "")
//}

// ServeHTTP gets called if a request receives the PEP. The function implements
// the PEP's main routine: It performs basic authentication, authorization with
// help of the PEP, transformation from SFCs into SFPs with help of the SFP
// Logic, and then forwards the package along the SFP.
func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "hello world")
	// Log all http requests incl. TLS informaion in the case of a successful TLS handshake
	//router.sysLogger.LogHTTPRequest(req)
	//
	//// Prepare Reponse
	//prepareResponse(w)
	//
	//// BLOCKING: Check if req.RemoteAddr is on one of the blocklists
	//if blocklist.BlockRequest(req) {
	//	io.WriteString(w, "Request has been rejected since you are on a blocklist. Contact your security advisor for more information.")
	//	w.WriteHeader(403)
	//	router.sysLogger.Infof("router: ServerHTTP(): Request from %s has been blocked since it is on a blocklist", req.RemoteAddr)
	//	return
	//}
	//
	//// Checks if the requested resource is even served by this PEP instance
	//serviceConf, ok := config.Config.ServiceSniMap[req.Host]
	//if !ok {
	//	io.WriteString(w, "Requested resource is not provided by this PEP. Please check again.")
	//	w.WriteHeader(404)
	//	router.sysLogger.Infof("router: ServerHTTP(): %s requested a resource that is not provided by this PEP.", req.RemoteAddr)
	//	return
	//}
	//
	//// Declares a new metadata instance which stores all necessary data to process the client's request
	//md := new(metadata.CpMetadata)
	//metadata.CollectMetadata(req, md)
	//
	//// AUTHENTICATION
	//// prompts the user for the authentication factors and evalautes them
	//// Check if the user is authenticated; if not authenticate her; if that fails return an error
	//// TODO: return error to client?
	//// Check if user has a valid session already
	//if !basic_auth.ClientHasValidSession(router.sysLogger, w, req, md) {
	//	//http.Redirect(w, req, "https://"+req.Host+"/9af1ecf7/login-portal", http.StatusFound) // 302
	//	//return
	//	basic_auth.PerformAuthentication(router.sysLogger, w, req, md)
	//	return
	//}
	//
	//// AUTHORIZATION
	//err := pdp.PerformAuthorization(router.sysLogger, req, md)
	//if err != nil {
	//	router.sysLogger.WithField("issuer", "PDP").Error(err)
	//	io.WriteString(w, "Internal ZTSFC error occured. Contact your security advisor if the error perists.")
	//	w.WriteHeader(500)
	//	return
	//}
	//
	//if !md.AuthDecision {
	//	// io.WriteString(w, md.AuthReason+". Contact your security advisor for more information.")
	//	router.sysLogger.Infof("router: ServeHTTP(): request from user %s was rejected due to the reason: %s", md.User, md.AuthReason)
	//	//w.WriteHeader(403)
	//	switch req.URL.Path {
	//	// Password Authentication
	//	case "/password-authentication":
	//		basic_auth.HandlePwdAuth(md.AuthReason, w)
	//		return
	//	// Passkey Authentication
	//	case "/passkey-authentication":
	//		basic_auth.HandlePasskeyAuthentication(md.AuthReason, w)
	//		return
	//	case "/begin-passkey-register":
	//		basic_auth.BeginPasskeyRegistration(w, req)
	//		return
	//	case "/finish-passkey-register":
	//		basic_auth.FinishPasskeyRegistration(w, req)
	//		return
	//	case "/begin-passkey-login":
	//		basic_auth.BeginPasskeyLogin(w, req)
	//		return
	//	case "/finish-passkey-login":
	//		basic_auth.FinishPasskeyLogin(router.sysLogger, w, req)
	//		return
	//	// All other cases for user without valid session
	//	default:
	//		basic_auth.HandleAuthenticationWelcome(md.AuthReason, w)
	//		return
	//	}
	//}
	//router.sysLogger.Debugf("router: ServeHTTP(): request from %s passed PDP. SFC: %s", req.RemoteAddr, md.SFC)
	//
	//// If user could be authenticated, define necessary variables for further processing
	//var proxy *httputil.ReverseProxy
	//var nextHopURL *url.URL
	//var certShownByPEPToNextHop tls.Certificate
	//var sfURLs []string
	//
	//// SFP LOGIC
	//// If SFC is empty, skip SFP Logic
	//if len(md.SFC) == 0 || config.Config.SfPool == nil {
	//	nextHopURL = serviceConf.TargetServiceUrl
	//	certShownByPEPToNextHop = serviceConf.X509KeyPairShownByPepToService
	//} else {
	//	err = sfpl.TransformSFCIntoSFP(router.sysLogger, md)
	//	if err != nil {
	//		router.sysLogger.Errorf("router: ServeHTTP(): could not transform SFC into SFP: %v", err)
	//		return
	//	}
	//	router.sysLogger.Debugf("router: ServeHTTP(): request passed SFP logic. SFP: %s", md.SFP)
	//
	//	if len(md.SFP) == 0 {
	//		router.sysLogger.Error("router: ServeHTTP(): SFP is empty altho SFC is not.")
	//		return
	//	}
	//
	//	// identify next hop, find its config and set nextHopURL and cert respectively
	//	nextHop := md.SFP[0]
	//	// TODO: make this dynamic later
	//	prepareSfMdHeader(req, md)
	//
	//	router.sysLogger.Debugf("router: ServeHTTP(): next hop: %s", nextHop)
	//
	//	nextHopConf, ok := config.Config.SfPool[nextHop.Name]
	//	if !ok {
	//		router.sysLogger.Errorf("next hop SF '%s' from the SFP does not exist in config file.", nextHop.Name)
	//		return
	//	}
	//
	//	nextHopURL, err = url.Parse(nextHop.URL)
	//	if err != nil {
	//		router.sysLogger.Errorf("router: ServeHTTP(): Could not parse URL '%s' value as URL: %v", nextHop.URL, err)
	//	}
	//
	//	certShownByPEPToNextHop = nextHopConf.X509KeyPairShownByPepToSf
	//
	//	// translate SF identifiers into ip addresses for remaining SFs
	//	for _, sf := range md.SFP[1:] {
	//		sfURLs = append(sfURLs, sf.URL)
	//	}
	//
	//	// finally append target service to list of SFP addresses, create a string of them and set this as header for following SFs
	//	sfURLs = append(sfURLs, serviceConf.TargetServiceAddr)
	//	addressesStr := strings.Join(sfURLs, ",")
	//	router.sysLogger.Debugf("router: ServeHTTP(): SFP as presented to following SFs: %s", addressesStr)
	//
	//	req.Header.Set("sfp", addressesStr)
	//
	//}
	//
	//router.sysLogger.Infof("router: ServeHTTP(): request from user %s (%s) forwarded to %s with SFP %v",
	//	md.User, md.Location, nextHopURL, sfURLs)
	//
	//proxy = httputil.NewSingleHostReverseProxy(nextHopURL)
	//
	//proxy.ErrorLog = log.New(router.sysLogger.GetWriter(), "", 0)
	//
	//// set proxy settings depending on the next hop scheme: http or https
	//// HTTPS
	//if nextHopURL.Scheme == "https" {
	//	// When the PEP is acting as a client; this defines his behavior
	//	proxy.Transport = &http.Transport{
	//		IdleConnTimeout:     10 * time.Second,
	//		MaxIdleConnsPerHost: 10000,
	//		TLSClientConfig: &tls.Config{
	//			Certificates:       []tls.Certificate{certShownByPEPToNextHop},
	//			InsecureSkipVerify: true,
	//			ClientAuth:         tls.RequireAndVerifyClientCert,
	//			ClientCAs:          config.Config.CAcertPoolPepAcceptsFromInt,
	//		},
	//	}
	//}
	//// HTTP
	//if nextHopURL.Scheme == "http" {
	//	proxy.Transport = &http.Transport{
	//		IdleConnTimeout:     10 * time.Second,
	//		MaxIdleConnsPerHost: 10000,
	//	}
	//}
	//
	//proxy.ServeHTTP(w, req)
}
