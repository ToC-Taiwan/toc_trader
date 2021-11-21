package main

import (
	"time"

	_ "gitlab.tocraw.com/root/toc_trader/init/sysinit"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/internal/network"

	"github.com/gin-gonic/gin"
	"gitlab.tocraw.com/root/toc_trader/init/sysparminit"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/routers"
)

func main() {
	// Gin server
	go func() {
		gin.SetMode(sysparminit.GlobalSettings.GetRunMode())
		g := gin.New()
		g.Use(routers.CorsMiddleware())
		g.Use(gin.Recovery())
		routers.AddSwagger(g)
		routers.InitRouters(g)
		if err := g.Run(":" + global.HTTPPort); err != nil {
			panic(err)
		}
	}()
	// Check Sinopac SRV is alive
	logger.GetLogger().Infof("Checking host on %s:%s...", global.PyServerHost, global.PyServerPort)
	for range time.Tick(time.Second) {
		if network.CheckPortIsOpen(global.PyServerHost, global.PyServerPort) {
			break
		}
	}
	// Main service
	go TradeProcess()
	// Keep thread running
	for {
		exit := <-global.ExitChannel
		if exit == global.ExitSignal {
			logger.GetLogger().Panic("manual exit")
		}
	}
}
