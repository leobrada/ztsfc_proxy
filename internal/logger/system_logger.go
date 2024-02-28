package logger

import (
	"fmt"
	"os"

	gct "github.com/leobrada/golang_convenience_tools"
	"github.com/sirupsen/logrus"
)

var (
	SystemLogger *logrus.Logger
	Success      string = "\033[32msuccessfully\033[0m"
)

func NewSystemLogger(systemLoggerOutput string, debugMode, errorMode bool, systemLoggerFormat string) error {
	SystemLogger = logrus.New()

	// set SystemLogger ouput path
	err := initSystemLoggerOutput(systemLoggerOutput)
	if err != nil {
		return fmt.Errorf("logger.InitSystemLogger(): %v", err)
	}

	// set SystemLogger info level
	err = initSystemLoggerInfoLevel(debugMode, errorMode)
	if err != nil {
		return fmt.Errorf("logger.InitSystemLogger(): %v", err)
	}

	// set SystemLogger output format
	initOutputFormat(systemLoggerFormat)

	SystemLogger.Debugf("logger.NewSystemLogger(): SystemLogger %s initialized. Output set to '%s', format set to '%s'", Success, systemLoggerOutput, systemLoggerFormat)

	return nil
}

func initSystemLoggerOutput(systemLoggerOutput string) error {
	loggerOutput, err := gct.GetWriter(systemLoggerOutput)
	if err != nil {
		SystemLogger.SetOutput(os.Stderr)
		return fmt.Errorf("logger.initSystemLoggerOutput(): %v", err)
	}
	SystemLogger.SetOutput(loggerOutput)
	return nil
}

func initSystemLoggerInfoLevel(debugMode, errorMode bool) error {
	switch {
	case debugMode && errorMode:
		return fmt.Errorf("logger.initSystemLoggerInfoLevel(): options 'debugMode' and 'errorMode' are mutually exclusive")
	case debugMode:
		SystemLogger.SetLevel(logrus.DebugLevel)
	case errorMode:
		SystemLogger.SetLevel(logrus.ErrorLevel)
	default:
		SystemLogger.SetLevel(logrus.InfoLevel)
	}
	return nil
}

func initOutputFormat(systemLoggerFormat string) {
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
		//SystemLogger.Debugln("logger.initOutputFormat(): SystemLogger output format set to 'json'")
	case "text":
		SystemLogger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
		//SystemLogger.Debugln("logger.initOutputFormat(): SystemLogger output format set to 'text'")
	default:
		SystemLogger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
		SystemLogger.Infof("logger.initOutputFormat(): SystemLogger output format '%s' is not supported. fallback to 'text'", systemLoggerFormat)
	}
}
