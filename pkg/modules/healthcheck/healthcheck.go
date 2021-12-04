// Package healthcheck package healthcheck
package healthcheck

import (
	"time"

	"gitlab.tocraw.com/root/toc_trader/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/sinopacapi"
)

var sinopacToken string

// ExitChannel ExitChannel
var ExitChannel chan int

func init() {
	ExitChannel = make(chan int)
}

// CheckSinopacToken CheckSinopacToken
func CheckSinopacToken() {
	tick := time.Tick(10 * time.Second)
	for range tick {
		token, err := sinopacapi.GetAgent().FetchServerKey()
		if err != nil {
			logger.GetLogger().Error(err)
			continue
		}
		if sinopacToken == "" {
			sinopacToken = token
		} else if sinopacToken != "" && sinopacToken != token {
			ExitService()
		}
	}
}

// ExitService ExitService
func ExitService() {
	ExitChannel <- global.ExitSignal
}
