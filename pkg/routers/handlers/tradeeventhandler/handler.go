// Package tradeeventhandler tradeeventhandler
package tradeeventhandler

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tocraw.com/root/toc_trader/pkg/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/tradeevent"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tradeeventprocess"
	"gitlab.tocraw.com/root/toc_trader/pkg/routers/handlers"
	"gitlab.tocraw.com/root/toc_trader/pkg/sinopacapi"
	"google.golang.org/protobuf/proto"
)

// AddHandlers AddHandlers
func AddHandlers(group *gin.RouterGroup) {
	group.POST("/trade-event", ReciveTradeEvent)
}

// ReciveTradeEvent ReciveTradeEvent
// @Summary ReciveTradeEvent
// @tags TradeEvent
// @accept json
// @produce json
// @param body body sinopacapi.EventProto true "Body"
// @success 200
// @failure 500 {object} handlers.ErrorResponse
// @Router /trade-event [post]
func ReciveTradeEvent(c *gin.Context) {
	req := sinopacapi.EventProto{}
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
	event := tradeevent.EventResponse{
		Event:     req.GetEvent(),
		EventCode: req.GetEventCode(),
		Info:      req.GetInfo(),
		Response:  req.GetRespCode(),
	}
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
		if err := sinopacapi.GetAgent().RestartSinopacSRV(); err != nil {
			logger.GetLogger().Panic(err)
		}
		logger.GetLogger().Error("Terminate, sinpac srv send 401, restart sinopac")
	}
	c.JSON(http.StatusOK, nil)
}
