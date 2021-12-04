// Package traderecordhandler package traderecordhandler
package traderecordhandler

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.tocraw.com/root/toc_trader/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/database"
	"gitlab.tocraw.com/root/toc_trader/pkg/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/traderecord"
	"gitlab.tocraw.com/root/toc_trader/pkg/routers/handlers"
	"gitlab.tocraw.com/root/toc_trader/pkg/sinopacapi"
	"google.golang.org/protobuf/proto"
)

// AddHandlers AddHandlers
func AddHandlers(group *gin.RouterGroup) {
	group.POST("/trade-record", UpdateTradeRecord)
}

// UpdateTradeRecord UpdateTradeRecord
// @Summary UpdateTradeRecord
// @tags TradeRecord
// @accept json
// @produce json
// @param body body sinopacapi.TradeRecordArrProto true "Body"
// @success 200
// @failure 500 {object} handlers.ErrorResponse
// @Router /trade-record [post]
func UpdateTradeRecord(c *gin.Context) {
	var res handlers.ErrorResponse
	body := sinopacapi.TradeRecordArrProto{}
	if byteArr, err := ioutil.ReadAll(c.Request.Body); err != nil {
		logger.GetLogger().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	} else if err := proto.Unmarshal(byteArr, &body); err != nil {
		logger.GetLogger().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	var records []*traderecord.TradeRecord
	for _, v := range body.Data {
		var orderTime time.Time
		name := global.AllStockNameMap.GetValueByKey(v.Code)
		orderTime, err := time.ParseInLocation(global.LongTimeLayout, v.OrderTime, time.Local)
		if err != nil {
			logger.GetLogger().Error(err)
			res.Response = err.Error()
			c.JSON(http.StatusInternalServerError, res)
			return
		}
		status := traderecord.StatusListMap[v.GetStatus()]
		action := traderecord.ActionListMap[v.GetAction()]
		tmp := traderecord.TradeRecord{
			StockNum:  v.Code,
			StockName: name,
			Action:    action,
			Price:     v.Price,
			Quantity:  v.Quantity,
			Status:    status,
			OrderID:   v.Id,
			OrderTime: orderTime,
		}
		records = append(records, &tmp)
	}
	if err := traderecord.UpdateByStockNumAndClose(records, database.GetAgent()); err != nil {
		logger.GetLogger().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	c.JSON(http.StatusOK, nil)
}
