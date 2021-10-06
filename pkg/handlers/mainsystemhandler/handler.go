// Package mainsystemhandler main handler
package mainsystemhandler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/handlers"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/process"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
)

// AddHandlers AddHandlers
func AddHandlers(group *gin.RouterGroup) {
	group.GET("/system/restart", Restart)
	group.GET("/system/trade/switch", GetTradeBotCondition)
	group.PUT("/system/trade/switch", UpdateTradeBotCondition)
	group.POST("/system/pyserver/host", UpdatePyServerHost)
}

// Restart Restart
// @Summary Restart
// @tags mainsystem
// @accept json
// @produce json
// @success 200
// @failure 500 {object} string
// @Router /system/restart [get]
func Restart(c *gin.Context) {
	var res handlers.ErrorResponse
	deployment := os.Getenv("DEPLOYMENT")
	if deployment != "docker" {
		res.Response = "you should be in the docker container"
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	process.RestartService()
	c.JSON(http.StatusOK, nil)
}

// UpdateTradeBotCondition UpdateTradeBotCondition
// @Summary UpdateTradeBotCondition
// @tags mainsystem
// @accept json
// @produce json
// @param body body UpdateTradeBotConditionBody true "Body"
// @success 200
// @failure 500 {object} handlers.ErrorResponse
// @Router /system/trade/switch [put]
func UpdateTradeBotCondition(c *gin.Context) {
	req := UpdateTradeBotConditionBody{}
	var res handlers.ErrorResponse
	if byteArr, err := ioutil.ReadAll(c.Request.Body); err != nil {
		logger.Logger.Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	} else if err := json.Unmarshal(byteArr, &req); err != nil {
		logger.Logger.Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	global.TradeSwitch.Buy = req.EnableBuy
	global.TradeSwitch.Sell = req.EnableSell
	global.TradeSwitch.UseBidAsk = req.UseBidAsk
	global.TradeSwitch.MeanTimeTradeStockNum = req.MeanTimeTradeStockNum
	logger.Logger.WithFields(map[string]interface{}{
		"EnableBuy":             global.TradeSwitch.Buy,
		"EnableSell":            global.TradeSwitch.Sell,
		"UseBidAsk":             global.TradeSwitch.UseBidAsk,
		"MeanTimeTradeStockNum": global.TradeSwitch.MeanTimeTradeStockNum,
	}).Info("Trade Switch Status")
	c.JSON(http.StatusOK, nil)
}

// GetTradeBotCondition GetTradeBotCondition
// @Summary GetTradeBotCondition
// @tags mainsystem
// @accept json
// @produce json
// @success 200
// @Router /system/trade/switch [get]
func GetTradeBotCondition(c *gin.Context) {
	data := UpdateTradeBotConditionBody{
		EnableBuy:             global.TradeSwitch.Buy,
		EnableSell:            global.TradeSwitch.Sell,
		UseBidAsk:             global.TradeSwitch.UseBidAsk,
		MeanTimeTradeStockNum: global.TradeSwitch.MeanTimeTradeStockNum,
	}
	c.JSON(http.StatusOK, data)
}

// UpdatePyServerHost UpdatePyServerHost
// @Summary UpdatePyServerHost
// @tags mainsystem
// @accept json
// @produce json
// @param py_host header string true "host"
// @success 200
// @failure 500 {object} handlers.ErrorResponse
// @Router /system/pyserver/host [post]
func UpdatePyServerHost(c *gin.Context) {
	var res handlers.ErrorResponse
	host := c.Request.Header.Get("py_host")
	logger.Logger.Warnf("Change PyServer to %s", host)
	if len(host) == 0 {
		res.Response = errors.New("host format is wrong").Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	global.PyServerHost = host
	c.JSON(http.StatusOK, nil)
}