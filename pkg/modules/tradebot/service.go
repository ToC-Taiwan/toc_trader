// Package tradebot package tradebot
package tradebot

import (
	"time"

	"github.com/markcheno/go-quote"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/analyzestreamtick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/streamtick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/traderecord"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tickanalyze"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
	"gitlab.tocraw.com/root/toc_trader/tools/stockutil"
)

// BuyOrderMap BuyOrderMap
var BuyOrderMap tradeRecordMutexStruct

// SellOrderMap SellOrderMap
var SellOrderMap tradeRecordMutexStruct

// FilledBuyOrderMap FilledBuyOrderMap
var FilledBuyOrderMap tradeRecordMutexStruct

// FilledSellOrderMap FilledSellOrderMap
var FilledSellOrderMap tradeRecordMutexStruct

// ManualSellMap ManualSellMap
var ManualSellMap tradeRecordMutexStruct

var lastTradeTime time.Time

// BuyBot BuyBot
func BuyBot(ch chan *analyzestreamtick.AnalyzeStreamTick) {
	for {
		analyzeTick := <-ch
		name := global.AllStockNameMap.GetName(analyzeTick.StockNum)
		outSum := analyzeTick.OutSum
		inSum := analyzeTick.InSum
		closeChangeRatio := analyzeTick.CloseChangeRatio
		outInRatio := analyzeTick.OutInRatio

		if !IsBuyPoint(analyzeTick, global.TickAnalyzeCondition) {
			continue
		}

		buyCost := GetStockBuyCost(analyzeTick.Close, global.OneTimeQuantity)
		if global.EnableBuy && !FilledBuyOrderMap.CheckStockExist(analyzeTick.StockNum) && BuyOrderMap.GetCount() < global.MeanTimeTradeStockNum &&
			!BuyOrderMap.CheckStockExist(analyzeTick.StockNum) && TradeQuota-buyCost > 0 {
			if order, err := PlaceOrder(BuyAction, analyzeTick.StockNum, global.OneTimeQuantity, analyzeTick.Close); err != nil {
				logger.Logger.WithFields(map[string]interface{}{
					"Msg":      err,
					"StockNum": analyzeTick.StockNum,
					"Quantity": global.OneTimeQuantity,
					"Price":    analyzeTick.Close,
				}).Error("Buy fail")
			} else if order.OrderID != "" && order.Status != traderecord.Failed {
				TradeQuota -= buyCost
				record := traderecord.TradeRecord{
					StockNum:  analyzeTick.StockNum,
					Price:     analyzeTick.Close,
					Quantity:  global.OneTimeQuantity,
					Action:    int64(BuyAction),
					BuyCost:   buyCost,
					TradeTime: time.Unix(0, analyzeTick.TimeStamp),
					OrderID:   order.OrderID,
				}
				BuyOrderMap.Set(record)
				go CheckBuyOrderStatus(record)
				tickTime := time.Unix(0, analyzeTick.TimeStamp).Local().Format(global.LongTimeLayout)
				replaceDate := tickTime[:10]
				clockTime := tickTime[11:19]
				logger.Logger.WithFields(map[string]interface{}{
					"Close":       analyzeTick.Close,
					"ChangeRatio": closeChangeRatio,
					"OutSum":      outSum,
					"InSum":       inSum,
					"OutInRatio":  outInRatio,
					"Name":        name,
				}).Infof("StreamTick Analyze: %s %s %s", replaceDate, clockTime, analyzeTick.StockNum)
				continue
			}
			logger.Logger.Warn("Buy Order is failed")
		}
	}
}

// IsBuyPoint IsBuyPoint
func IsBuyPoint(analyzeTick *analyzestreamtick.AnalyzeStreamTick, cond global.AnalyzeCondition) bool {
	closeChangeRatio := analyzeTick.CloseChangeRatio
	// if global.UseBidAsk {
	// 	lastBidAsk := bidaskprocess.TmpBidAskMap.GetLastOneByStockNum(analyzeTick.StockNum)
	// 	if !lastBidAsk.IsBestBid() {
	// 		logger.Logger.Warn("Not best bid")
	// 		return false
	// 	}
	// }
	if analyzeTick.OpenChangeRatio > cond.OpenChangeRatio || closeChangeRatio < cond.CloseChangeRatioLow || closeChangeRatio > cond.CloseChangeRatioHigh {
		return false
	}
	if analyzeTick.OutInRatio < cond.OutInRatio || analyzeTick.OutSum < cond.OutSum || analyzeTick.CloseDiff <= cond.CloseDiff {
		return false
	}
	if analyzeTick.Rsi > float64(cond.RsiLow) {
		return false
	}
	return true
}

// SellBot SellBot
func SellBot(ch chan *streamtick.StreamTick) {
	var historyClose []float64
	for {
		tick := <-ch
		historyClose = append(historyClose, tick.Close)
		if len(historyClose) > global.HistoryCloseCount {
			historyClose = historyClose[1:]
		}
		filled, err := traderecord.CheckIsFilledByOrderID(BuyOrderMap.GetOrderID(tick.StockNum), global.GlobalDB)
		if err != nil {
			logger.Logger.Error(err)
			continue
		}
		if filled && !SellOrderMap.CheckStockExist(tick.StockNum) && global.EnableSell {
			originalOrderClose := BuyOrderMap.GetClose(tick.StockNum)
			sellPrice := GetSellPrice(tick, BuyOrderMap.GetTradeTime(tick.StockNum), historyClose, originalOrderClose, global.TickAnalyzeCondition)
			if sellPrice == 0 {
				continue
			} else if order, err := PlaceOrder(SellAction, tick.StockNum, global.OneTimeQuantity, sellPrice); err != nil {
				logger.Logger.WithFields(map[string]interface{}{
					"Msg":      err,
					"Stock":    tick.StockNum,
					"Quantity": global.OneTimeQuantity,
					"Price":    sellPrice,
				}).Error("Sell fail")
			} else if order.OrderID != "" && order.Status != traderecord.Failed {
				record := traderecord.TradeRecord{
					StockNum:  tick.StockNum,
					Price:     sellPrice,
					Quantity:  global.OneTimeQuantity,
					Action:    int64(SellAction),
					TradeTime: time.Unix(0, tick.TimeStamp),
					OrderID:   order.OrderID,
				}
				SellOrderMap.Set(record)
				go CheckSellOrderStatus(record)
				continue
			}
			logger.Logger.Warn("Sell Order is failed")
		}
	}
}

// GetSellPrice GetSellPrice
func GetSellPrice(tick *streamtick.StreamTick, tradeTime time.Time, historyClose []float64, originalOrderClose float64, cond global.AnalyzeCondition) float64 {
	tickTimeUnix := time.Unix(0, tick.TimeStamp)
	lastTime := time.Date(tickTimeUnix.Year(), tickTimeUnix.Month(), tickTimeUnix.Day(), 13, 0, 0, 0, time.Local)
	if len(historyClose) < global.HistoryCloseCount && tickTimeUnix.Before(lastTime) {
		return 0
	}
	var sellPrice float64
	var input quote.Quote
	input.Close = historyClose
	switch {
	case tick.Close < stockutil.GetNewClose(originalOrderClose, -1) && tickanalyze.GenerateRSI(input) > float64(cond.RsiHigh):
		sellPrice = tick.Close
	case tickanalyze.GenerateRSI(input) > float64(cond.RsiHigh):
		sellPrice = tick.Close
	case tradeTime.Add(600*time.Second).Before(tickTimeUnix) && tick.Close >= stockutil.GetNewClose(originalOrderClose, 2):
		sellPrice = tick.Close
	case tradeTime.Add(1800*time.Second).Before(tickTimeUnix) && tick.Close >= stockutil.GetNewClose(originalOrderClose, 1):
		sellPrice = tick.Close
	case ManualSellMap.CheckStockExist(tick.StockNum):
		sellPrice = ManualSellMap.GetClose(tick.StockNum)
		if sellPrice == 0 {
			sellPrice = tick.Close
		}
	case tickTimeUnix.After(lastTime):
		sellPrice = tick.Close
	default:
		sellPrice = 0
	}
	return sellPrice
}

// CheckBuyOrderStatus CheckBuyOrderStatus
func CheckBuyOrderStatus(record traderecord.TradeRecord) {
	for {
		time.Sleep(1 * time.Second)
		order, err := traderecord.GetOrderByOrderID(record.OrderID, global.GlobalDB)
		if err != nil {
			logger.Logger.Error(err)
			continue
		}
		if order.Status == 4 {
			TradeQuota += record.BuyCost
			BuyOrderMap.Delete(record.StockNum)
			logger.Logger.WithFields(map[string]interface{}{
				"StockNum": order.StockNum,
				"Name":     order.StockName,
				"Quantity": order.Quantity,
				"Price":    order.Price,
			}).Info("Place Order Fail")
			return
		}
		if record.TradeTime.Add(60*time.Second).Before(time.Now()) && order.Status != 6 && order.Status != 5 {
			if err := Cancel(record.OrderID); err != nil {
				logger.Logger.Error(err)
				return
			}
			TradeQuota += record.BuyCost
			BuyOrderMap.Delete(record.StockNum)
			logger.Logger.WithFields(map[string]interface{}{
				"StockNum": order.StockNum,
				"Name":     order.StockName,
				"Quantity": order.Quantity,
				"Price":    order.Price,
			}).Info("Cancel Buy Order Success")
			return
		}
		if order.Status == 6 {
			FilledBuyOrderMap.Set(order)
			logger.Logger.WithFields(map[string]interface{}{
				"StockNum": order.StockNum,
				"Name":     order.StockName,
				"Quantity": order.Quantity,
				"Price":    order.Price,
			}).Info("Buy Stock Success")
			return
		}
	}
}

// CheckSellOrderStatus CheckSellOrderStatus
func CheckSellOrderStatus(record traderecord.TradeRecord) {
	for {
		time.Sleep(1 * time.Second)
		order, err := traderecord.GetOrderByOrderID(record.OrderID, global.GlobalDB)
		if err != nil {
			logger.Logger.Error(err)
			continue
		}
		if order.Status == 4 {
			TradeQuota += record.BuyCost
			BuyOrderMap.Delete(record.StockNum)
			logger.Logger.WithFields(map[string]interface{}{
				"StockNum": order.StockNum,
				"Name":     order.StockName,
				"Quantity": order.Quantity,
				"Price":    order.Price,
			}).Info("Place Order Fail")
			return
		}
		if record.TradeTime.Add(60*time.Second).Before(time.Now()) && order.Status != 6 && order.Status != 5 {
			if err := Cancel(record.OrderID); err != nil {
				logger.Logger.Error(err)
				return
			}
			SellOrderMap.Delete(record.StockNum)
			logger.Logger.WithFields(map[string]interface{}{
				"StockNum": order.StockNum,
				"Name":     order.StockName,
				"Quantity": order.Quantity,
				"Price":    order.Price,
			}).Info("Cancel Sell Order Success")
			return
		}
		if order.Status == 6 {
			FilledSellOrderMap.Set(order)
			BuyOrderMap.Delete(record.StockNum)
			SellOrderMap.Delete(record.StockNum)
			if ManualSellMap.CheckStockExist(record.StockNum) {
				ManualSellMap.Delete(record.StockNum)
			}
			logger.Logger.WithFields(map[string]interface{}{
				"StockNum": order.StockNum,
				"Name":     order.StockName,
				"Quantity": order.Quantity,
				"Price":    order.Price,
			}).Info("Sell Stock Success")
			return
		}
	}
}
