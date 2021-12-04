// Package manualtradehandler package manualtradehandler
package manualtradehandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tocraw.com/root/toc_trader/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/traderecord"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tradebot"
	"gitlab.tocraw.com/root/toc_trader/pkg/routers/handlers"
)

// AddHandlers AddHandlers
func AddHandlers(group *gin.RouterGroup) {
	group.POST("/manual/sell", ManualSellStock)
	group.POST("/manual/buy-later", ManualBuyLaterStock)
}

// ManualSellStock ManualSellStock
// @Summary ManualSellStock
// @tags ManualTrade
// @accept json
// @produce json
// @param body body ManualSellBody true "Body"
// @success 200
// @failure 500 {object} handlers.ErrorResponse
// @Router /manual/sell [post]
func ManualSellStock(c *gin.Context) {
	req := ManualSellBody{}
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

	record := traderecord.TradeRecord{
		StockNum: req.StockNum,
		Price:    req.Price,
	}
	tradebot.ManualSellMap.Set(record)
	logger.GetLogger().WithFields(map[string]interface{}{
		"StockNum": req.StockNum,
		"Price":    req.Price,
		"Name":     global.AllStockNameMap.GetValueByKey(req.StockNum),
	}).Info("Manual Sell")
	c.JSON(http.StatusOK, nil)
}

// ManualBuyLaterStock ManualBuyLaterStock
// @Summary ManualBuyLaterStock
// @tags ManualTrade
// @accept json
// @produce json
// @param body body ManualBuyLaterBody true "Body"
// @success 200
// @failure 500 {object} handlers.ErrorResponse
// @Router /manual/buy-later [post]
func ManualBuyLaterStock(c *gin.Context) {
	req := ManualBuyLaterBody{}
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

	record := traderecord.TradeRecord{
		StockNum: req.StockNum,
		Price:    req.Price,
	}
	tradebot.ManualBuyLaterMap.Set(record)
	logger.GetLogger().WithFields(map[string]interface{}{
		"StockNum": req.StockNum,
		"Price":    req.Price,
		"Name":     global.AllStockNameMap.GetValueByKey(req.StockNum),
	}).Info("Manual Buy Later")
	c.JSON(http.StatusOK, nil)
}
