// Package tradebothandler tradebothandler
package tradebothandler

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/handlers"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/bidask"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/streamtick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/bidaskprocess"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/subscribe"
	"google.golang.org/protobuf/proto"
)

// AddHandlers AddHandlers
func AddHandlers(group *gin.RouterGroup) {
	group.POST("/data/streamtick", ReceiveStreamTick)
	group.POST("/data/bid-ask", ReceiveBidAsk)
}

// ReceiveStreamTick ReceiveStreamTick
// @Summary ReceiveStreamTick
// @tags Data
// @accept json
// @produce json
// @param body body streamtick.StreamTickProto true "Body"
// @success 200
// @failure 500 {object} handlers.ErrorResponse
// @Router /data/streamtick [post]
func ReceiveStreamTick(c *gin.Context) {
	req := streamtick.StreamTickProto{}
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
	tmp, err := req.ProtoToStreamTick()
	if err != nil {
		logger.GetLogger().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	if tmp.Simtrade == 1 {
		subscribe.SimTradeChannel <- 1
		return
	}
	if tmp.TimeStamp != 0 {
		if forwardCh := *subscribe.ForwardStreamTickChannelMap.GetChannelByStockNum(req.Tick.Code); forwardCh != nil {
			forwardCh <- tmp
		}
		if reverseCh := *subscribe.ReverseStreamTickChannelMap.GetChannelByStockNum(req.Tick.Code); reverseCh != nil {
			reverseCh <- tmp
		}
	}
	c.JSON(http.StatusOK, nil)
}

// ReceiveBidAsk ReceiveBidAsk
// @Summary ReceiveBidAsk
// @tags Data
// @accept json
// @produce json
// @param body body bidask.BidAskProto true "Body"
// @success 200
// @failure 500 {object} handlers.ErrorResponse
// @Router /data/bid-ask [post]
func ReceiveBidAsk(c *gin.Context) {
	req := bidask.BidAskProto{}
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
	if req.BidAsk.Simtrade == 1 {
		return
	}
	if len(req.BidAsk.BidPrice) != 5 || len(req.BidAsk.BidVolume) != 5 || len(req.BidAsk.DiffBidVol) != 5 ||
		len(req.BidAsk.AskPrice) != 5 || len(req.BidAsk.AskVolume) != 5 || len(req.BidAsk.DiffAskVol) != 5 {
		logger.GetLogger().Error("Data is broken")
		return
	}
	data, err := req.BidAsk.ToBidAsk()
	if err != nil {
		logger.GetLogger().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	bidaskprocess.TmpBidAskMap.Append(data)
	c.JSON(http.StatusOK, nil)
}
