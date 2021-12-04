package main

import (
	"os"
	"time"

	_ "gitlab.tocraw.com/root/toc_trader/init/sysinit"

	"gitlab.tocraw.com/root/toc_trader/global"
	"gitlab.tocraw.com/root/toc_trader/init/sysparminit"
	"gitlab.tocraw.com/root/toc_trader/pkg/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/healthcheck"
	"gitlab.tocraw.com/root/toc_trader/pkg/restclient"
	"gitlab.tocraw.com/root/toc_trader/pkg/routers"
	"gitlab.tocraw.com/root/toc_trader/pkg/sinopacapi"
	"gitlab.tocraw.com/root/toc_trader/pkg/utils"
)

func main() {
	// gin server
	routers.Init()
	// Get sinopac server host and port from sysparm
	serverHost := sysparminit.GlobalSettings.GetPyServerHost()
	serverPort := sysparminit.GlobalSettings.GetPyServerPort()
	// Check Sinopac SRV is alive
	logger.GetLogger().Infof("Checking host on %s:%s...", serverHost, serverPort)
	for range time.Tick(time.Second) {
		if utils.CheckPortIsOpen(serverHost, serverPort) {
			logger.GetLogger().Info("Sinopac SRV is alive")
			break
		}
	}
	// create new sinopac agent
	sinopacapi.NewAgent(serverHost, serverPort, restclient.GetClient())
	// send ip to sinopac srv
	if err := sinopacapi.GetAgent().UpdateTraderIP(utils.GetHostIP()); err != nil {
		logger.GetLogger().Panic(err)
	}
	// Check token is expired or not, if expired, restart service
	go healthcheck.CheckSinopacToken()
	// Main service
	go TradeProcess()

	// Keep thread running
	for {
		signal := <-healthcheck.ExitChannel
		if signal == global.ExitSignal {
			logger.GetLogger().Warn("Manual exit")
			os.Exit(1)
		}
	}
}
