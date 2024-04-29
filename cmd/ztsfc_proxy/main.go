package main

import (
	"flag"

	"github.com/leobrada/ztsfc_proxy/internal/configs"
	"github.com/leobrada/ztsfc_proxy/internal/frontend"
	"github.com/leobrada/ztsfc_proxy/internal/logger"
)

var (
	// variables storing command-line arguments
	confFilePath       string
	systemLoggerOutput string
	debugMode          bool
	errorMode          bool
	systemLoggerFormat string

	// pointer to central configuration file storing all configuration settings
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

	// loads all configuration settings from the yaml file specified in 'confFilePath'
	config, err = configs.NewConfig(confFilePath)
	if err != nil {
		logger.SystemLogger.Fatalf("main.init(): %v", err)
	}

	logger.SystemLogger.Infof("main.init(): Configuration %s initialized from from %s - OK", logger.Success, confFilePath)
}

func main() {
	// starts frontend HTTP server with all necessary configuration settings set like TLS config
	// derived from config file read in from the 'confFilePath' variable in main.init()
	frontend, err := frontend.NewFrontend(config)
	if err != nil {
		logger.SystemLogger.Fatalf("main.main(): %v", err)
	}

	if err = frontend.ListenAndServeTLS("", ""); err != nil {
		logger.SystemLogger.Fatalf("main.main(): %v", err)
	}
}
