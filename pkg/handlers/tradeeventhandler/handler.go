// Package tradeeventhandler tradeeventhandler
package tradeeventhandler

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tocraw.com/root/toc_trader/internal/healthcheck"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/handlers"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/tradeevent"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tradeeventprocess"
	"google.golang.org/protobuf/proto"
)

// AddHandlers AddHandlers
func AddHandlers(group *gin.RouterGroup) {
	group.POST("/trade-event", ReciveTradeEvent)
}

// ReciveTradeEvent ReciveTradeEvent
// @Summary ReciveTradeEvent
// @tags tradeevent
// @accept json
// @produce json
// @param body body tradeevent.EventProto true "Body"
// @success 200
// @failure 500 {object} handlers.ErrorResponse
// @Router /trade-event [post]
func ReciveTradeEvent(c *gin.Context) {
	req := tradeevent.EventProto{}
	var res handlers.ErrorResponse
	if byteArr, err := ioutil.ReadAll(c.Request.Body); err != nil {
		logger.GetLogger().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	} else if err := proto.Unmarshal(byteArr, &req); err != nil {
		logger.GetLogger().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	event := req.ToEventResponse()
	if err := tradeeventprocess.TradeEventSaver(event); err != nil {
		logger.GetLogger().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	logger.GetLogger().WithFields(map[string]interface{}{
		"Resopnse":  req.RespCode,
		"Info":      req.Info,
		"EventCode": req.EventCode,
		"Event":     req.Event,
	}).Info("SinoPacSRV Event")
	if req.EventCode == 401 {
		if err := healthcheck.AskSinopacSRVRestart(); err != nil {
			logger.GetLogger().Panic(err)
		}
		logger.GetLogger().Error("Terminate, sinpac srv send 401, restart sinopac")
	}
	c.JSON(http.StatusOK, nil)
}
