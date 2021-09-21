// Package logger is logger
package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// Logger Global logger
var Logger = log.New()

func init() {
	Logger.SetFormatter(&log.TextFormatter{
		TimestampFormat:  "2006/01/02 15:04:05",
		FullTimestamp:    true,
		QuoteEmptyFields: true,
		DisableColors:    false,
		PadLevelText:     false,
	})
	Logger.SetLevel(log.TraceLevel)
	Logger.SetOutput(os.Stdout)
	// deployment := os.Getenv("DEPLOYMENT")
	// var caPath string
	// switch {
	// case deployment == "docker":
	// 	caPath = "/toc_trader/configs/certs/ca.crt"
	// case deployment == "test":
	// 	caPath = "/builds/root/toc_trader/configs/certs/ca.crt"
	// default:
	// 	caPath = "/Users/timhsu/dev_projects/golang/toc_trader/configs/certs/ca.crt"
	// }
	// p := logrusmqtt.MQTTHookParams{
	// 	Hostname:   sysparminit.GlobalSettings.GetPyServerHost(),
	// 	Port:       8887,
	// 	Username:   "tradebotsino",
	// 	Password:   "asdf0000",
	// 	Topic:      "runtime",
	// 	CAFilepath: caPath,
	// 	Insecure:   true,
	// 	QoS:        0,
	// 	Retain:     false,
	// }

	// hook, err := logrusmqtt.NewMQTTHook(p, log.TraceLevel)
	// if err != nil {
	// 	panic(err)
	// }
	// Logger.AddHook(hook)
}
