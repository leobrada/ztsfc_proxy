package ztsfc_logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func logger_init() {
	Log = logrus.New()
	Log.Out = os.Stdout
	Log.Level = logrus.InfoLevel
	Log.Formatter = &logrus.JSONFormatter{}
}
