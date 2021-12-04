// Package tradebalancehandler package tradebalancehandler
package tradebalancehandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"gitlab.tocraw.com/root/toc_trader/pkg/database"
	"gitlab.tocraw.com/root/toc_trader/pkg/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/balance"
	"gitlab.tocraw.com/root/toc_trader/pkg/routers/handlers"
)

// AddHandlers AddHandlers
func AddHandlers(group *gin.RouterGroup) {
	group.GET("/balance", GetAllBalance)
	group.POST("/balance", ImportBalance)
	group.DELETE("/balance", DeletaAllBalance)
}

// GetAllBalance GetAllBalance
// @Summary GetAllBalance
// @tags Balance
// @accept json
// @produce json
// @success 200 {object} []balance.Balance
// @failure 500 {object} handlers.ErrorResponse
// @Router /balance [get]
func GetAllBalance(c *gin.Context) {
	var res handlers.ErrorResponse
	allBalance, err := balance.GetAllBalance(database.GetAgent())
	if err != nil {
		logger.GetLogger().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	if len(allBalance) > 1 {
		sort.Slice(allBalance, func(i, j int) bool {
			return allBalance[i].TradeDay.Before(allBalance[j].TradeDay)
		})
	}
	c.JSON(http.StatusOK, allBalance)
}

// ImportBalance ImportBalance
// @Summary ImportBalance
// @tags Balance
// @accept json
// @produce json
// @param body body []balance.Balance{} true "Body"
// @success 200
// @failure 500 {object} handlers.ErrorResponse
// @Router /balance [post]
func ImportBalance(c *gin.Context) {
	var res handlers.ErrorResponse
	body := []balance.Balance{}
	if byteArr, err := ioutil.ReadAll(c.Request.Body); err != nil {
		logger.GetLogger().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	} else if err := json.Unmarshal(byteArr, &body); err != nil {
		logger.GetLogger().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	if len(body) > 1 {
		sort.Slice(body, func(i, j int) bool {
			return body[i].TradeDay.Before(body[j].TradeDay)
		})
	}
	for _, v := range body {
		if err := balance.InsertOrUpdate(v, database.GetAgent()); err != nil {
			logger.GetLogger().Error(err)
			res.Response = err.Error()
			c.JSON(http.StatusInternalServerError, res)
			return
		}
	}
	c.JSON(http.StatusOK, nil)
}

// DeletaAllBalance DeletaAllBalance
// @Summary DeletaAllBalance
// @tags Balance
// @accept json
// @produce json
// @success 200
// @failure 500 {object} handlers.ErrorResponse
// @Router /balance [delete]
func DeletaAllBalance(c *gin.Context) {
	var res handlers.ErrorResponse
	if err := balance.DeleteAll(database.GetAgent()); err != nil {
		logger.GetLogger().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	c.JSON(http.StatusOK, nil)
}
