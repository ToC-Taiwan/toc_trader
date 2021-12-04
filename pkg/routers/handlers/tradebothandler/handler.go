// Package tradebothandler tradebothandler
package tradebothandler

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tocraw.com/root/toc_trader/pkg/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/bidask"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/streamtick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/bidaskprocess"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/subscribe"
	"gitlab.tocraw.com/root/toc_trader/pkg/routers/handlers"
	"gitlab.tocraw.com/root/toc_trader/pkg/sinopacapi"
	"gitlab.tocraw.com/root/toc_trader/pkg/utils"
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
// @param body body sinopacapi.StreamTickProto true "Body"
// @success 200
// @failure 500 {object} handlers.ErrorResponse
// @Router /data/streamtick [post]
func ReceiveStreamTick(c *gin.Context) {
	req := sinopacapi.StreamTickProto{}
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
	if req.Tick.GetSimtrade() == 1 {
		subscribe.SimTradeChannel <- 1
		return
	}
	ts, err := utils.MicroDateTimeToTimeStamp(req.Tick.GetDateTime())
	if err != nil {
		logger.GetLogger().Error(err)
		return
	}
	tmp := streamtick.StreamTick{
		StockNum:        req.Tick.GetCode(),
		TimeStamp:       ts,
		Open:            req.Tick.GetOpen(),
		AvgPrice:        req.Tick.GetAvgPrice(),
		Close:           req.Tick.GetClose(),
		High:            req.Tick.GetHigh(),
		Low:             req.Tick.GetLow(),
		Amount:          req.Tick.GetAmount(),
		AmountSum:       req.Tick.GetTotalAmount(),
		Volume:          req.Tick.GetVolume(),
		VolumeSum:       req.Tick.GetTotalVolume(),
		TickType:        req.Tick.GetTickType(),
		ChgType:         req.Tick.GetChgType(),
		PriceChg:        req.Tick.GetPriceChg(),
		PctChg:          req.Tick.GetPctChg(),
		BidSideTotalVol: req.Tick.GetBidSideTotalVol(),
		AskSideTotalVol: req.Tick.GetAskSideTotalVol(),
		BidSideTotalCnt: req.Tick.GetBidSideTotalCnt(),
		AskSideTotalCnt: req.Tick.GetAskSideTotalCnt(),
		Suspend:         req.Tick.GetSuspend(),
		Simtrade:        req.Tick.GetSimtrade(),
	}
	if tmp.TimeStamp != 0 {
		if forwardCh := *subscribe.ForwardStreamTickChannelMap.GetChannelByStockNum(req.Tick.Code); forwardCh != nil {
			forwardCh <- &tmp
		}
		if reverseCh := *subscribe.ReverseStreamTickChannelMap.GetChannelByStockNum(req.Tick.Code); reverseCh != nil {
			reverseCh <- &tmp
		}
	}
	c.JSON(http.StatusOK, nil)
}

// ReceiveBidAsk ReceiveBidAsk
// @Summary ReceiveBidAsk
// @tags Data
// @accept json
// @produce json
// @param body body sinopacapi.BidAskProto true "Body"
// @success 200
// @failure 500 {object} handlers.ErrorResponse
// @Router /data/bid-ask [post]
func ReceiveBidAsk(c *gin.Context) {
	req := sinopacapi.BidAskProto{}
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
	if req.BidAsk.GetSimtrade() == 1 {
		return
	}
	if len(req.BidAsk.BidPrice) != 5 || len(req.BidAsk.BidVolume) != 5 || len(req.BidAsk.DiffBidVol) != 5 ||
		len(req.BidAsk.AskPrice) != 5 || len(req.BidAsk.AskVolume) != 5 || len(req.BidAsk.DiffAskVol) != 5 {
		logger.GetLogger().Error("Data is broken")
		return
	}
	ts, err := utils.MicroDateTimeToTimeStamp(req.BidAsk.GetDateTime())
	if err != nil {
		logger.GetLogger().Error(err)
		return
	}
	data := bidask.BidAsk{
		StockNum:    req.BidAsk.Code,
		TimeStamp:   ts,
		BidPrice1:   req.BidAsk.GetBidPrice()[4],
		BidVolume1:  req.BidAsk.GetBidVolume()[4],
		DiffBidVol1: req.BidAsk.GetDiffBidVol()[4],

		BidPrice2:   req.BidAsk.GetBidPrice()[3],
		BidVolume2:  req.BidAsk.GetBidVolume()[3],
		DiffBidVol2: req.BidAsk.GetDiffBidVol()[3],

		BidPrice3:   req.BidAsk.GetBidPrice()[2],
		BidVolume3:  req.BidAsk.GetBidVolume()[2],
		DiffBidVol3: req.BidAsk.GetDiffBidVol()[2],

		BidPrice4:   req.BidAsk.GetBidPrice()[1],
		BidVolume4:  req.BidAsk.GetBidVolume()[1],
		DiffBidVol4: req.BidAsk.GetDiffBidVol()[1],

		BidPrice5:   req.BidAsk.GetBidPrice()[0],
		BidVolume5:  req.BidAsk.GetBidVolume()[0],
		DiffBidVol5: req.BidAsk.GetDiffBidVol()[0],

		AskPrice1:   req.BidAsk.GetAskPrice()[0],
		AskVolume1:  req.BidAsk.GetAskVolume()[0],
		DiffAskVol1: req.BidAsk.GetDiffAskVol()[0],

		AskPrice2:   req.BidAsk.GetAskPrice()[1],
		AskVolume2:  req.BidAsk.GetAskVolume()[1],
		DiffAskVol2: req.BidAsk.GetDiffAskVol()[1],

		AskPrice3:   req.BidAsk.GetAskPrice()[2],
		AskVolume3:  req.BidAsk.GetAskVolume()[2],
		DiffAskVol3: req.BidAsk.GetDiffAskVol()[2],

		AskPrice4:   req.BidAsk.GetAskPrice()[3],
		AskVolume4:  req.BidAsk.GetAskVolume()[3],
		DiffAskVol4: req.BidAsk.GetDiffAskVol()[3],

		AskPrice5:   req.BidAsk.GetAskPrice()[4],
		AskVolume5:  req.BidAsk.GetAskVolume()[4],
		DiffAskVol5: req.BidAsk.GetDiffAskVol()[4],

		Suspend:  req.BidAsk.GetSuspend(),
		Simtrade: req.BidAsk.GetSimtrade(),
	}
	bidaskprocess.TmpBidAskMap.Append(&data)
	c.JSON(http.StatusOK, nil)
}
