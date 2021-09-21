// Package routers package routers
package routers

import (
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"gitlab.tocraw.com/root/toc_trader/docs"
	"gitlab.tocraw.com/root/toc_trader/init/sysparminit"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
)

// AddSwagger AddSwagger
// @title ToC Trader
// @version 0.1.0
// @description API docs for ToC Trader
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /trade-bot
func AddSwagger(router *gin.Engine) {
	deployment := os.Getenv("DEPLOYMENT")
	docs.SwaggerInfo.Host = "toc-trader.tocraw.com:" + sysparminit.GlobalSettings.GetHTTPPort()
	if deployment != "docker" {
		docs.SwaggerInfo.Host = "localhost:" + sysparminit.GlobalSettings.GetHTTPPort()
	}
	url := ginSwagger.URL("/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	logger.Logger.Info("http://" + docs.SwaggerInfo.Host + "/swagger/index.html")
}
