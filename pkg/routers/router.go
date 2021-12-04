// Package routers all routers
package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tocraw.com/root/toc_trader/init/sysparminit"
	"gitlab.tocraw.com/root/toc_trader/pkg/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/routers/handlers/mainsystemhandler"
	"gitlab.tocraw.com/root/toc_trader/pkg/routers/handlers/manualtradehandler"
	"gitlab.tocraw.com/root/toc_trader/pkg/routers/handlers/targethandler"
	"gitlab.tocraw.com/root/toc_trader/pkg/routers/handlers/tradebalancehandler"
	"gitlab.tocraw.com/root/toc_trader/pkg/routers/handlers/tradebothandler"
	"gitlab.tocraw.com/root/toc_trader/pkg/routers/handlers/tradecondhandler"
	"gitlab.tocraw.com/root/toc_trader/pkg/routers/handlers/tradeeventhandler"
	"gitlab.tocraw.com/root/toc_trader/pkg/routers/handlers/traderecordhandler"
)

// Init Init
func Init() {
	go func() {
		gin.SetMode(sysparminit.GlobalSettings.GetRunMode())
		g := gin.New()
		g.Use(corsMiddleware())
		g.Use(gin.Recovery())
		err := g.SetTrustedProxies(nil)
		if err != nil {
			logger.GetLogger().Panic(err)
		}
		addSwagger(g)
		initRouters(g)
		if err := g.Run(":" + sysparminit.GlobalSettings.GetHTTPPort()); err != nil {
			logger.GetLogger().Panic(err)
		}
	}()
}

func initRouters(router *gin.Engine) {
	mainRoute := router.Group("trade-bot")
	mainsystemhandler.AddHandlers(mainRoute)
	tradebothandler.AddHandlers(mainRoute)
	tradeeventhandler.AddHandlers(mainRoute)
	traderecordhandler.AddHandlers(mainRoute)
	tradebalancehandler.AddHandlers(mainRoute)
	tradecondhandler.AddHandlers(mainRoute)
	manualtradehandler.AddHandlers(mainRoute)
	targethandler.AddHandlers(mainRoute)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
		c.Set("content-type", "application/json")
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		c.Next()
	}
}
