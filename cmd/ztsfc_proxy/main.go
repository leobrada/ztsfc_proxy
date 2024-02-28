package main

import (
	"flag"

	"github.com/leobrada/ztsfc_proxy/internal/configs"
	"github.com/leobrada/ztsfc_proxy/internal/logger"
	"github.com/leobrada/ztsfc_proxy/internal/router"
)

var (
	confFilePath       string
	systemLoggerOutput string
	debugMode          bool
	errorMode          bool
	systemLoggerFormat string

	// config pointer
	config *configs.Config
)

func init() {
	// parse command-line arguments
	flag.StringVar(&confFilePath, "c", "./configs/config.yml", "Path to user defined YML config file")
	flag.StringVar(&systemLoggerOutput, "o", "stdout", "Output path of system logger")
	flag.BoolVar(&debugMode, "d", false, "Enable debug output level")
	flag.BoolVar(&errorMode, "e", false, "Enable error output level")
	flag.StringVar(&systemLoggerFormat, "f", "text", "Output format of system logger {text,json}")
	flag.Parse()

	// initialize the global SystemLogger defined in package "logger"
	err := logger.NewSystemLogger(systemLoggerOutput, debugMode, errorMode, systemLoggerFormat)
	if err != nil {
		logger.SystemLogger.Fatalf("main.init(): %v", err)
	}

	config, err = configs.NewConfig(confFilePath)
	if err != nil {
		logger.SystemLogger.Fatalf("main.init(): %v", err)
	}

	logger.SystemLogger.Infof("main.init(): Configuration %s initialized from from %s - OK", logger.Success, confFilePath)
}

func main() {
	ztsfcProxy := router.NewRouter()
	logger.SystemLogger.Info("main.main(): New ZTSFC Proxy successfully created")
	if err := ztsfcProxy.ListenAndServeTLS(); err != nil {
		logger.SystemLogger.Errorf("main.main(): Error on listening and server: %v", err)
	}
}
