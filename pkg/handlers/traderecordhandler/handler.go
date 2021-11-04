// Package traderecordhandler package traderecordhandler
package traderecordhandler

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tocraw.com/root/toc_trader/internal/database"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/handlers"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/traderecord"
	"google.golang.org/protobuf/proto"
)

// AddHandlers AddHandlers
func AddHandlers(group *gin.RouterGroup) {
	group.POST("/trade-record", UpdateTradeRecord)
}

// UpdateTradeRecord UpdateTradeRecord
// @Summary UpdateTradeRecord
// @tags traderecord
// @accept json
// @produce json
// @param body body traderecord.TradeRecordArrProto true "Body"
// @success 200
// @failure 500 {object} handlers.ErrorResponse
// @Router /trade-record [post]
func UpdateTradeRecord(c *gin.Context) {
	var res handlers.ErrorResponse
	body := traderecord.TradeRecordArrProto{}
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
	records, err := body.ToTradeRecordFromProto()
	if err != nil {
		logger.GetLogger().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	if err := traderecord.UpdateByStockNumAndClose(records, database.GetAgent()); err != nil {
		logger.GetLogger().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	c.JSON(http.StatusOK, nil)
}
