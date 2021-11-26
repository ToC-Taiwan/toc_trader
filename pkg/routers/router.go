// Package routers all routers
package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tocraw.com/root/toc_trader/pkg/handlers/mainsystemhandler"
	"gitlab.tocraw.com/root/toc_trader/pkg/handlers/manualtradehandler"
	"gitlab.tocraw.com/root/toc_trader/pkg/handlers/targethandler"
	"gitlab.tocraw.com/root/toc_trader/pkg/handlers/tradebalancehandler"
	"gitlab.tocraw.com/root/toc_trader/pkg/handlers/tradebothandler"
	"gitlab.tocraw.com/root/toc_trader/pkg/handlers/tradecondhandler"
	"gitlab.tocraw.com/root/toc_trader/pkg/handlers/tradeeventhandler"
	"gitlab.tocraw.com/root/toc_trader/pkg/handlers/traderecordhandler"
)

// InitRouters InitRouters
func InitRouters(router *gin.Engine) {
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

// CorsMiddleware CorsMiddleware
func CorsMiddleware() gin.HandlerFunc {
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
