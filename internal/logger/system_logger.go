package logger

import (
	"fmt"
	"os"

	gct "github.com/leobrada/golang_convenience_tools"
	"github.com/sirupsen/logrus"
)

var SystemLogger *logrus.Logger

func SystemLoggerInit(systemLoggerOutput string, debugMode bool, systemLoggerFormat string) error {
	SystemLogger = logrus.New()

	// set SystemLogger ouput path
	loggerOutput, err := gct.GetWriter(systemLoggerOutput)
	if err != nil {
		SystemLogger.SetOutput(os.Stderr)
		return fmt.Errorf("logger.System_logger_init(): %v", err)
	}
	SystemLogger.SetOutput(loggerOutput)

	// set SystemLogger info level
	if debugMode {
		SystemLogger.SetLevel(logrus.DebugLevel)
	} else {
		SystemLogger.SetLevel(logrus.InfoLevel)
	}

	// set SystemLogger output format
	switch systemLoggerFormat {
	case "json":
		SystemLogger.SetFormatter(&logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "@timestamp",
				logrus.FieldKeyLevel: "@level",
				logrus.FieldKeyMsg:   "@message",
				logrus.FieldKeyFunc:  "@caller",
			},
		})
		SystemLogger.Infoln("logger.System_logger_init(): SystemLogger output format set to 'json'")
	case "text":
		SystemLogger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
		SystemLogger.Infoln("logger.System_logger_init(): SystemLogger output format set to 'text'")
	default:
		SystemLogger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
		SystemLogger.Infof("logger.System_logger_init(): SystemLogger output format '%s' is not supported. fallback to 'text'", systemLoggerFormat)
	}

	return nil
}
