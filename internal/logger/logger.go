// Package logger is logger
package logger

import (
	"os"
	"path/filepath"

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
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	basePath := filepath.Dir(ex)
	Log = logrus.New()
	deployment := os.Getenv("DEPLOYMENT")
	if deployment == "docker" {
		Log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: global.LongTimeLayout,
			PrettyPrint:     false,
		})
	} else {
		Log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat:  "2006/01/02 15:04:05",
			FullTimestamp:    true,
			QuoteEmptyFields: true,
			PadLevelText:     false,
			ForceColors:      true,
			ForceQuote:       true,
			// CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			// 	fileName := path.Base(frame.File)
			// 	return frame.Function, fileName
			// },
		})
	}
	// Log.SetReportCaller(true)
	Log.SetLevel(logrus.TraceLevel)
	Log.SetOutput(os.Stdout)
	pathMap := lfshook.PathMap{
		logrus.PanicLevel: basePath + "/logs/panic.json",
		logrus.FatalLevel: basePath + "/logs/fetal.json",
		logrus.ErrorLevel: basePath + "/logs/error.json",
		logrus.WarnLevel:  basePath + "/logs/warn.json",
		logrus.InfoLevel:  basePath + "/logs/info.json",
		logrus.DebugLevel: basePath + "/logs/debug.json",
		logrus.TraceLevel: basePath + "/logs/error.json",
	}
	Log.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	))
	return Log
}
