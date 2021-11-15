// Package tradebothandler tradebothandler
package tradebothandler

import (
	"crypto/rand"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tocraw.com/root/toc_trader/internal/common"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/handlers"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/bidask"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/streamtick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/traderecord"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/bidaskprocess"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/subscribe"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tradebot"
	"google.golang.org/protobuf/proto"
)

// AddHandlers AddHandlers
func AddHandlers(group *gin.RouterGroup) {
	group.POST("/data/streamtick", ReceiveStreamTick)
	group.GET("/data/target", GetTarget)
	group.POST("/data/bid-ask", ReceiveBidAsk)

	group.POST("/manual/sell", ManualSellStock)
	group.POST("/manual/buy-later", ManualBuyLaterStock)

	group.GET("/switch", GetTradeBotSwitch)
	group.PUT("/switch", UpdateTradeBotCondition)

	group.GET("/condition", GetTradeCondition)
}

// ReceiveStreamTick ReceiveStreamTick
// @Summary ReceiveStreamTick
// @tags tradebot
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
		*subscribe.ForwardStreamTickChannelMap.GetChannelByStockNum(req.Tick.Code) <- tmp
		*subscribe.ReverseStreamTickChannelMap.GetChannelByStockNum(req.Tick.Code) <- tmp
	}
	c.JSON(http.StatusOK, nil)
}

// GetTarget GetTarget
// @Summary GetTarget
// @tags tradebot
// @accept json
// @produce json
// @param count header string true "count"
// @success 200 {object} []TargetResponse
// @failure 500 {object} handlers.ErrorResponse
// @Router /data/target [get]
func GetTarget(c *gin.Context) {
	var res handlers.ErrorResponse
	count := c.Request.Header.Get("count")
	countInt64, err := common.StrToInt64(count)
	if err != nil {
		logger.GetLogger().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	total := len(global.TargetArr)
	response := []TargetResponse{}
	already := make(map[int64]bool)
	if int(countInt64) >= total {
		countInt64 = int64(total)
	}
	for i := 0; i < int(countInt64); i++ {
		randomBigInt, err := rand.Int(rand.Reader, big.NewInt((int64(total))))
		if err != nil {
			logger.GetLogger().Error(err)
			res.Response = err.Error()
			c.JSON(http.StatusInternalServerError, res)
			return
		}
		random := randomBigInt.Int64()
		if _, ok := already[random]; ok {
			i--
			continue
		}
		data := TargetResponse{
			StockNum: global.TargetArr[random],
			Close:    global.StockCloseByDateMap.GetClose(global.TargetArr[random], global.LastTradeDay.Format(global.ShortTimeLayout)),
		}
		already[random] = true
		response = append(response, data)
	}
	c.JSON(http.StatusOK, response)
}

// ReceiveBidAsk ReceiveBidAsk
// @Summary ReceiveBidAsk
// @tags tradebot
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

// ManualSellStock ManualSellStock
// @Summary ManualSellStock
// @tags tradebot
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
		"Name":     global.AllStockNameMap.GetName(req.StockNum),
	}).Info("Manual Sell")
	c.JSON(http.StatusOK, nil)
}

// ManualBuyLaterStock ManualBuyLaterStock
// @Summary ManualBuyLaterStock
// @tags tradebot
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
		"Name":     global.AllStockNameMap.GetName(req.StockNum),
	}).Info("Manual Buy Later")
	c.JSON(http.StatusOK, nil)
}

// GetTradeBotSwitch GetTradeBotSwitch
// @Summary GetTradeBotSwitch
// @tags tradebot
// @accept json
// @produce json
// @success 200
// @Router /switch [get]
func GetTradeBotSwitch(c *gin.Context) {
	data := UpdateTradeBotSwitchBody{
		EnableBuy:                    global.TradeSwitch.Buy,
		EnableSell:                   global.TradeSwitch.Sell,
		EnableSellFirst:              global.TradeSwitch.SellFirst,
		EnableBuyLater:               global.TradeSwitch.BuyLater,
		UseBidAsk:                    global.TradeSwitch.UseBidAsk,
		MeanTimeTradeStockNum:        global.TradeSwitch.MeanTimeTradeStockNum,
		MeanTimeReverseTradeStockNum: global.TradeSwitch.MeanTimeReverseTradeStockNum,
	}
	c.JSON(http.StatusOK, data)
}

// UpdateTradeBotCondition UpdateTradeBotCondition
// @Summary UpdateTradeBotCondition
// @tags tradebot
// @accept json
// @produce json
// @param body body UpdateTradeBotSwitchBody true "Body"
// @success 200
// @failure 500 {object} handlers.ErrorResponse
// @Router /switch [put]
func UpdateTradeBotCondition(c *gin.Context) {
	req := UpdateTradeBotSwitchBody{}
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
	global.TradeSwitch.Buy = req.EnableBuy
	global.TradeSwitch.Sell = req.EnableSell
	global.TradeSwitch.SellFirst = req.EnableSell
	global.TradeSwitch.BuyLater = req.EnableBuyLater
	global.TradeSwitch.UseBidAsk = req.UseBidAsk
	global.TradeSwitch.MeanTimeTradeStockNum = req.MeanTimeTradeStockNum
	global.TradeSwitch.MeanTimeReverseTradeStockNum = req.MeanTimeReverseTradeStockNum
	logger.GetLogger().WithFields(map[string]interface{}{
		"EnableBuy":             global.TradeSwitch.Buy,
		"EnableSell":            global.TradeSwitch.Sell,
		"EnableSellFirst":       global.TradeSwitch.SellFirst,
		"EnableBuyLater":        global.TradeSwitch.BuyLater,
		"MeanTimeTradeStockNum": global.TradeSwitch.MeanTimeTradeStockNum,
	}).Info("Trade Switch Status")
	c.JSON(http.StatusOK, nil)
}

// GetTradeCondition GetTradeCondition
// @Summary GetTradeCondition
// @tags tradebot
// @accept json
// @produce json
// @success 200
// @Router /condition [get]
func GetTradeCondition(c *gin.Context) {
	data := []simulationcond.AnalyzeCondition{
		global.ForwardCond,
		global.ReverseCond,
	}
	c.JSON(http.StatusOK, data)
}
