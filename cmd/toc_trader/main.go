package main

import (
	"os"
	"time"

	_ "gitlab.tocraw.com/root/toc_trader/init/sysinit"
	"gitlab.tocraw.com/root/toc_trader/internal/healthcheck"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/internal/network"

	"github.com/gin-gonic/gin"
	"gitlab.tocraw.com/root/toc_trader/init/sysparminit"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/routers"
)

var sinopacToken string

func main() {
	// Gin server
	go func() {
		gin.SetMode(sysparminit.GlobalSettings.GetRunMode())
		g := gin.New()
		g.Use(routers.CorsMiddleware())
		g.Use(gin.Recovery())
		err := g.SetTrustedProxies(nil)
		if err != nil {
			logger.GetLogger().Panic(err)
		}
		routers.AddSwagger(g)
		routers.InitRouters(g)
		if err := g.Run(":" + global.HTTPPort); err != nil {
			logger.GetLogger().Panic(err)
		}
	}()
	// Check Sinopac SRV is alive
	logger.GetLogger().Infof("Checking host on %s:%s...", global.PyServerHost, global.PyServerPort)
	for range time.Tick(time.Second) {
		if network.CheckPortIsOpen(global.PyServerHost, global.PyServerPort) {
			break
		}
	}
	// Send ip to sinopac srv
	if err := sendHostIP(getHostIP()); err != nil {
		logger.GetLogger().Panic(err)
	}
	// Check token is expired or not, if expired, restart service
	go checkSinopacToken()
	// Main service
	go TradeProcess()
	// Keep thread running
	for {
		signal := <-global.ExitChannel
		if signal == global.ExitSignal {
			logger.GetLogger().Warn("Manual exit")
			os.Exit(1)
		}
	}
}

func checkSinopacToken() {
	for range time.Tick(10 * time.Second) {
		token, err := healthcheck.GetSinopacSRVToken()
		if err != nil {
			logger.GetLogger().Panic(err)
		}
		if sinopacToken == "" {
			sinopacToken = token
		} else if sinopacToken != "" && sinopacToken != token {
			healthcheck.ExitService()
		}
	}
}
