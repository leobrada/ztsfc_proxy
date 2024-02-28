package logger

import (
	"fmt"
	"log"
	"os"

	gct "github.com/leobrada/golang_convenience_tools"
	"github.com/leobrada/ztsfc_proxy/internal/configs"
)

func NewDataPlaneLogger(config *configs.Config) (*log.Logger, error) {
	loggerOutput, err := gct.GetWriter(config.DataPlaneLogger.Output)
	if err != nil {
		SystemLogger.SetOutput(os.Stderr)
		return nil, fmt.Errorf("logger.NewDataPlaneLogger(): %v", err)
	}

	dpLogger := log.New(loggerOutput, "\033[31m[DataPlane]\033[0m - ", 0b11100011)

	SystemLogger.Debugf("logger.NewSystemLogger(): DataPaneLogger %s initialized. Output set to '%s'", Success, config.DataPlaneLogger.Output)

	return dpLogger, nil
}
