package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func Logger_init() {
	Log = logrus.New()
	Log.Out = os.Stdout
	Log.Level = logrus.InfoLevel
	Log.Formatter = &logrus.JSONFormatter{}
}
