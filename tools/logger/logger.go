// Package logger is logger
package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
)

// Logger Global logger
var Logger = log.New()

func init() {
	deployment := os.Getenv("DEPLOYMENT")
	if deployment == "docker" {
		Logger.SetFormatter(&log.JSONFormatter{
			TimestampFormat: global.LongTimeLayout,
			PrettyPrint:     false,
		})
	} else {
		Logger.SetFormatter(&log.TextFormatter{
			TimestampFormat:  "2006/01/02 15:04:05",
			FullTimestamp:    true,
			QuoteEmptyFields: true,
			PadLevelText:     false,
		})
	}
	Logger.SetLevel(log.TraceLevel)
	// Logger.SetOutput(os.Stdout)
}
