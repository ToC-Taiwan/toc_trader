// Package logger is logger
package logger

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"gitlab.tocraw.com/root/toc_trader/global"
)

var client *logrus.Logger

// GetLogger GetLogger
func GetLogger() *logrus.Logger {
	if client != nil {
		return client
	}
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	basePath := filepath.Dir(ex)
	client = logrus.New()
	deployment := os.Getenv("DEPLOYMENT")
	if deployment == "docker" {
		client.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: global.LongTimeLayout,
			PrettyPrint:     false,
		})
	} else {
		client.SetFormatter(&logrus.TextFormatter{
			TimestampFormat:  "2006/01/02 15:04:05",
			FullTimestamp:    true,
			QuoteEmptyFields: true,
			PadLevelText:     false,
			ForceColors:      true,
			ForceQuote:       true,
		})
	}
	// Log.SetReportCaller(true)
	fileNamePrefix := time.Now().Format(time.RFC3339)[:16] + "-"
	fileNamePrefix = strings.ReplaceAll(fileNamePrefix, ":", "")
	client.SetLevel(logrus.TraceLevel)
	client.SetOutput(os.Stdout)
	pathMap := lfshook.PathMap{
		logrus.PanicLevel: basePath + "/logs/" + fileNamePrefix + "panic.json",
		logrus.FatalLevel: basePath + "/logs/" + fileNamePrefix + "fetal.json",
		logrus.ErrorLevel: basePath + "/logs/" + fileNamePrefix + "error.json",
		logrus.WarnLevel:  basePath + "/logs/" + fileNamePrefix + "warn.json",
		logrus.InfoLevel:  basePath + "/logs/" + fileNamePrefix + "info.json",
		logrus.DebugLevel: basePath + "/logs/" + fileNamePrefix + "debug.json",
		logrus.TraceLevel: basePath + "/logs/" + fileNamePrefix + "error.json",
	}
	client.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	))
	return client
}
