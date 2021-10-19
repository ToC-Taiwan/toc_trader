// Package logger is logger
package logger

import (
	"os"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
)

// Log Log
var Log *logrus.Logger

// GetLogger GetLogger
func GetLogger() *logrus.Logger {
	if Log != nil {
		return Log
	}
	var basePath string
	Log = logrus.New()
	deployment := os.Getenv("DEPLOYMENT")
	if deployment == "docker" {
		basePath = "/toc_trader"
		Log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: global.LongTimeLayout,
			PrettyPrint:     false,
		})
	} else {
		basePath = "./"
		Log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat:  "2006/01/02 15:04:05",
			FullTimestamp:    true,
			QuoteEmptyFields: true,
			PadLevelText:     false,
			ForceColors:      true,
		})
	}
	Log.SetLevel(logrus.TraceLevel)
	Log.SetOutput(os.Stdout)
	pathMap := lfshook.PathMap{
		logrus.PanicLevel: basePath + "/logs/panic.log",
		logrus.FatalLevel: basePath + "/logs/fetal.log",
		logrus.ErrorLevel: basePath + "/logs/error.log",
		logrus.WarnLevel:  basePath + "/logs/warn.log",
		logrus.InfoLevel:  basePath + "/logs/info.log",
		logrus.DebugLevel: basePath + "/logs/debug.log",
		logrus.TraceLevel: basePath + "/logs/error.log",
	}
	Log.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	))
	return Log
}
