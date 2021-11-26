// Package mainsystemhandler main handler
package mainsystemhandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gitlab.tocraw.com/root/toc_trader/internal/healthcheck"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/handlers"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/sysparm"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/parameters"
)

// AddHandlers AddHandlers
func AddHandlers(group *gin.RouterGroup) {
	group.GET("/system/restart", Restart)
	group.POST("/system/sysparm", UpdateSysparm)

	group.GET("/trade/switch", GetTradeBotSwitch)
	group.PUT("/trade/switch", UpdateTradeBotSwitch)
}

// Restart Restart
// @Summary Restart
// @tags MainSystem
// @accept json
// @produce json
// @success 200
// @failure 500 {object} string
// @Router /system/restart [get]
func Restart(c *gin.Context) {
	var res handlers.ErrorResponse
	deployment := os.Getenv("DEPLOYMENT")
	if deployment != "docker" {
		res.Response = "you should be in the docker container(from toc_trader)"
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	healthcheck.ExitService()
	c.JSON(http.StatusOK, nil)
}

// UpdateSysparm UpdateSysparm
// @Summary UpdateSysparm
// @tags MainSystem
// @accept json
// @produce json
// @param body body []sysparm.Parameters true "Body"
// @success 200
// @failure 500 {object} handlers.ErrorResponse
// @Router /system/sysparm [post]
func UpdateSysparm(c *gin.Context) {
	var req []sysparm.Parameters
	var res handlers.ErrorResponse
	if byteArr, err := ioutil.ReadAll(c.Request.Body); err != nil {
		logger.GetLogger().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	} else if err := json.Unmarshal(byteArr, &req); err != nil {
		logger.GetLogger().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	for _, v := range req {
		if err := parameters.UpdateSysparm(v.Key, v.Value); err != nil {
			logger.GetLogger().Error(err)
			res.Response = err.Error()
			c.JSON(http.StatusInternalServerError, res)
			return
		}
	}
	c.JSON(http.StatusOK, nil)
}

// GetTradeBotSwitch GetTradeBotSwitch
// @Summary GetTradeBotSwitch
// @tags MainSystem
// @accept json
// @produce json
// @success 200 {object} global.SystemSwitch
// @Router /trade/switch [get]
func GetTradeBotSwitch(c *gin.Context) {
	c.JSON(http.StatusOK, global.TradeSwitch)
}

// UpdateTradeBotSwitch UpdateTradeBotSwitch
// @Summary UpdateTradeBotSwitch
// @tags MainSystem
// @accept json
// @produce json
// @param body body global.SystemSwitch true "Body"
// @success 200
// @failure 500 {object} handlers.ErrorResponse
// @Router /trade/switch [put]
func UpdateTradeBotSwitch(c *gin.Context) {
	req := global.SystemSwitch{}
	var res handlers.ErrorResponse
	if byteArr, err := ioutil.ReadAll(c.Request.Body); err != nil {
		logger.GetLogger().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	} else if err := json.Unmarshal(byteArr, &req); err != nil {
		logger.GetLogger().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	global.TradeSwitch = req
	logger.GetLogger().WithFields(map[string]interface{}{
		"EnableBuy":             global.TradeSwitch.Buy,
		"EnableSell":            global.TradeSwitch.Sell,
		"EnableSellFirst":       global.TradeSwitch.SellFirst,
		"EnableBuyLater":        global.TradeSwitch.BuyLater,
		"MeanTimeTradeStockNum": global.TradeSwitch.MeanTimeTradeStockNum,
	}).Info("Trade Switch Status")
	c.JSON(http.StatusOK, nil)
}
