// Package tradecondhandler package tradecondhandler
package tradecondhandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"gitlab.tocraw.com/root/toc_trader/internal/database"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/handlers"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulate"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
)

// AddHandlers AddHandlers
func AddHandlers(group *gin.RouterGroup) {
	group.GET("/condition/latest", GetLatestTradeCondition)

	group.GET("/condition", GetAllBestTradeCondition)
	group.POST("/condition", ImpoprtTradeCondition)
	group.DELETE("/condition", DeletaAllResultAndCond)
}

// GetLatestTradeCondition GetLatestTradeCondition
// @Summary GetLatestTradeCondition
// @tags condition
// @accept json
// @produce json
// @success 200 {object} []simulationcond.AnalyzeCondition
// @Router /condition/latest [get]
func GetLatestTradeCondition(c *gin.Context) {
	data := []simulationcond.AnalyzeCondition{
		global.ForwardCond,
		global.ReverseCond,
	}
	c.JSON(http.StatusOK, data)
}

// GetAllBestTradeCondition GetAllBestTradeCondition
// @Summary GetLatestTradeCondition
// @tags condition
// @accept json
// @produce json
// @success 200 {object} []simulate.Result
// @Router /condition [get]
func GetAllBestTradeCondition(c *gin.Context) {
	var res handlers.ErrorResponse
	var data []simulate.Result
	forwardResultCond, err := simulate.GetAllBestForwardResultAndCond(database.GetAgent())
	if err != nil {
		logger.GetLogger().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	reverseResultCond, err := simulate.GetAllBestReverseResultAndCond(database.GetAgent())
	if err != nil {
		logger.GetLogger().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	data = append(data, forwardResultCond...)
	data = append(data, reverseResultCond...)
	c.JSON(http.StatusOK, data)
}

// ImpoprtTradeCondition ImpoprtTradeCondition
// @Summary ImpoprtTradeCondition
// @tags condition
// @accept json
// @produce json
// @param body body []simulate.Result{} true "Body"
// @success 200
// @Router /condition [post]
func ImpoprtTradeCondition(c *gin.Context) {
	var res handlers.ErrorResponse
	body := []simulate.Result{}
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
		tmp := v
		if err := simulationcond.Insert(&tmp.Cond, database.GetAgent()); err != nil {
			logger.GetLogger().Error(err)
			res.Response = err.Error()
			c.JSON(http.StatusInternalServerError, res)
			return
		}
		if err := simulate.Insert(&tmp, database.GetAgent()); err != nil {
			logger.GetLogger().Error(err)
			res.Response = err.Error()
			c.JSON(http.StatusInternalServerError, res)
			return
		}
	}
	c.JSON(http.StatusOK, nil)
}

// DeletaAllResultAndCond DeletaAllResultAndCond
// @Summary DeletaAllResultAndCond
// @tags condition
// @accept json
// @produce json
// @success 200
// @failure 500 {object} handlers.ErrorResponse
// @Router /condition [delete]
func DeletaAllResultAndCond(c *gin.Context) {
	var res handlers.ErrorResponse
	if err := simulate.DeleteAll(database.GetAgent()); err != nil {
		logger.GetLogger().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	if err := simulationcond.DeleteAll(database.GetAgent()); err != nil {
		logger.GetLogger().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	c.JSON(http.StatusOK, nil)
}
