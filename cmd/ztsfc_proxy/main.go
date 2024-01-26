package main

import (
	"flag"
	"os"

	"github.com/leobrada/ztsfc_proxy/config"
	"github.com/leobrada/ztsfc_proxy/internal/logger"
	"github.com/leobrada/ztsfc_proxy/internal/router"
)

var (
	confFilePath       string
	systemLoggerOutput string
	debugMode          bool
	systemLoggerFormat string
)

func init() {
	// parse command-line arguments
	flag.StringVar(&confFilePath, "config", "./config/config.yml", "Path to user defined YML config file")
	flag.StringVar(&confFilePath, "c", "./config/config.yml", "Path to user defined YML config file")
	flag.StringVar(&systemLoggerOutput, "output", "stdout", "Output path of system logger")
	flag.StringVar(&systemLoggerOutput, "o", "stdout", "Output path of system logger")
	flag.BoolVar(&debugMode, "debug", false, "Enable debug output level")
	flag.BoolVar(&debugMode, "d", false, "Enable debug output level")
	flag.StringVar(&systemLoggerFormat, "format", "text", "Output format of system logger {text,json}")
	flag.StringVar(&systemLoggerFormat, "f", "text", "Output format of system logger {text,json}")
	flag.Parse()

	// initialize the global SystemLogger defined in package "logger"
	err := logger.SystemLoggerInit(systemLoggerOutput, debugMode, systemLoggerFormat)
	if err != nil {
		logger.SystemLogger.Errorf("main.init(): %v", err)
		os.Exit(1)
	}

	err = config.Config.InitGlobalConfig(confFilePath)
	if err != nil {
		logger.SystemLogger.Errorf("main.init(): %v", err)
		os.Exit(1)
	}

	logger.SystemLogger.Infof("main.init(): Initializing ZTSFC Proxy from %s - OK", confFilePath)
}

func main() {
	ztsfcProxy := router.NewRouter()
	logger.SystemLogger.Info("main.main(): New ZTSFC Proxy successfully created")
	if err := ztsfcProxy.ListenAndServe(); err != nil {
		logger.SystemLogger.Errorf("main.main(): Error on listening and server: %v", err)
	}
}

/*
func main() {

	server := &http.Server{
		Addr:     "134.60.77.40:8081",
		ErrorLog: log.New(logger.Out, "", 0),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Log the request
		logger.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Info("Request received")

		// Respond with "hello world"
		fmt.Fprint(w, "hello world")
	})

	if err := server.ListenAndServe(); err != nil {
		logger.Error(err)
	}
}
*/
