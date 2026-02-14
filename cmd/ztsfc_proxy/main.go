// ztsfc_proxy is the main entrypoint for the proxy binary.
// It wires configuration, logging, and the frontend server together.
package main

import (
	"flag"

	"github.com/leobrada/ztsfc_proxy/internal/configs"
	"github.com/leobrada/ztsfc_proxy/internal/frontend"
	"github.com/leobrada/ztsfc_proxy/internal/logger"
)

var (
	// CLI flags for runtime configuration.
	confFilePath       string
	systemLoggerOutput string
	debugMode          bool
	errorMode          bool
	systemLoggerFormat string

	// Parsed configuration loaded from the YAML file.
	config *configs.Config
)

func init() {
	// Parse command-line arguments before any initialization happens.
	flag.StringVar(&confFilePath, "c", "./configs/config.yml", "Path to user defined YML config file")
	flag.StringVar(&systemLoggerOutput, "o", "stdout", "Output path of system logger")
	flag.BoolVar(&debugMode, "d", false, "Enable debug output level")
	flag.BoolVar(&errorMode, "e", false, "Enable error output level")
	flag.StringVar(&systemLoggerFormat, "f", "text", "Output format of system logger {text,json}")
	flag.Parse()

	// Initialize the global system logger used by the application.
	err := logger.NewSystemLogger(systemLoggerOutput, debugMode, errorMode, systemLoggerFormat)
	if err != nil {
		logger.SystemLogger.Fatalf("main.init(): %v", err)
	}

	// Load all configuration settings from the YAML file.
	config, err = configs.NewConfig(confFilePath)
	if err != nil {
		logger.SystemLogger.Fatalf("main.init(): %v", err)
	}

	logger.SystemLogger.Infof("main.init(): Configuration %s initialized from from %s - OK", logger.Success, confFilePath)
}

func main() {
	// Start the frontend HTTP server using the loaded configuration.
	// TLS settings are derived from the config initialized in init().
	frontend, err := frontend.NewFrontend(config)
	if err != nil {
		logger.SystemLogger.Fatalf("main.main(): %v", err)
	}

	// Empty cert/key args here because TLS config is already set on the server.
	if err = frontend.ListenAndServeTLS("", ""); err != nil {
		logger.SystemLogger.Fatalf("main.main(): %v", err)
	}
}
