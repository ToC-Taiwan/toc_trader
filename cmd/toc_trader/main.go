package main

import (
	_ "gitlab.tocraw.com/root/toc_trader/init/sysinit"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
	"gitlab.tocraw.com/root/toc_trader/tools/network"

	"time"

	"github.com/gin-gonic/gin"
	"gitlab.tocraw.com/root/toc_trader/init/sysparminit"
	"gitlab.tocraw.com/root/toc_trader/pkg/core"
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
	for range time.NewTicker(time.Second).C {
		if network.CheckPortIsOpen(global.PyServerHost, global.PyServerPort) {
			break
		}
	}
	// Main service
	go core.TradeProcess()
	// Keep thread running
	for {
		exit := <-global.ExitChannel
		if exit == global.ExitSignal {
			logger.Logger.Panic("manual exit")
		}
	}
}
