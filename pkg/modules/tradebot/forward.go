// Package tradebot package tradebot
package tradebot

import (
	"time"

	"github.com/markcheno/go-quote"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/analyzestreamtick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/streamtick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/traderecord"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tickanalyze"
	"gitlab.tocraw.com/root/toc_trader/tools/db"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
	"gitlab.tocraw.com/root/toc_trader/tools/stockutil"
)

// BuyOrderMap BuyOrderMap
var BuyOrderMap tradeRecordMutexMap

// SellOrderMap SellOrderMap
var SellOrderMap tradeRecordMutexMap

// FilledBuyOrderMap FilledBuyOrderMap
var FilledBuyOrderMap tradeRecordMutexMap

// FilledSellOrderMap FilledSellOrderMap
var FilledSellOrderMap tradeRecordMutexMap

// ManualSellMap ManualSellMap
var ManualSellMap tradeRecordMutexMap

// IsBuyPoint IsBuyPoint
func IsBuyPoint(analyzeTick *analyzestreamtick.AnalyzeStreamTick, cond simulationcond.AnalyzeCondition) bool {
	closeChangeRatio := analyzeTick.CloseChangeRatio
	if analyzeTick.Volume < cond.VolumePerSecond*int64(analyzeTick.TotalTime) {
		return false
	}
	if analyzeTick.OpenChangeRatio > cond.OpenChangeRatio || closeChangeRatio < cond.CloseChangeRatioLow || closeChangeRatio > cond.CloseChangeRatioHigh {
		return false
	}
	if analyzeTick.OutInRatio < cond.OutInRatio || analyzeTick.CloseDiff <= cond.CloseDiff {
		return false
	}
	if analyzeTick.Rsi > cond.RsiLow {
		return false
	}
	return true
}

// BuyBot BuyBot
func BuyBot(analyzeTick *analyzestreamtick.AnalyzeStreamTick) {
	buyCost := GetStockBuyCost(analyzeTick.Close, global.OneTimeQuantity)
	if BuyOrderMap.GetCount() < global.TradeSwitch.MeanTimeTradeStockNum && TradeQuota-buyCost > 0 {
		if order, err := PlaceOrder(BuyAction, analyzeTick.StockNum, global.OneTimeQuantity, analyzeTick.Close); err != nil {
			logger.GetLogger().WithFields(map[string]interface{}{
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
		}
	} else {
		logger.GetLogger().Warn("Over MeanTimeTradeStockNum or Quota")
	}
}

// SellBot SellBot
func SellBot(ch chan *streamtick.StreamTick, cond simulationcond.AnalyzeCondition) {
	var historyClose []float64
	for {
		tick := <-ch
		historyClose = append(historyClose, tick.Close)
		if len(historyClose) > int(cond.HistoryCloseCount) {
			historyClose = historyClose[1:]
		}
		filled, err := traderecord.CheckIsFilledByOrderID(BuyOrderMap.GetOrderIDByStockNum(tick.StockNum), db.GetAgent())
		if err != nil {
			logger.GetLogger().Error(err)
			continue
		}
		if filled && !SellOrderMap.CheckStockExist(tick.StockNum) {
			originalOrderClose := BuyOrderMap.GetClose(tick.StockNum)
			sellPrice := GetSellPrice(tick, BuyOrderMap.GetTradeTime(tick.StockNum), historyClose, originalOrderClose, cond)
			if sellPrice == 0 {
				continue
			} else if order, err := PlaceOrder(SellAction, tick.StockNum, global.OneTimeQuantity, sellPrice); err != nil {
				logger.GetLogger().WithFields(map[string]interface{}{
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
			}
		}
	}
}

// GetSellPrice GetSellPrice
func GetSellPrice(tick *streamtick.StreamTick, tradeTime time.Time, historyClose []float64, originalOrderClose float64, cond simulationcond.AnalyzeCondition) float64 {
	tickTimeUnix := time.Unix(0, tick.TimeStamp)
	lastTime := time.Date(tickTimeUnix.Year(), tickTimeUnix.Month(), tickTimeUnix.Day(), global.TradeOutEndHour, global.TradeOutEndMinute, 0, 0, time.Local)
	if len(historyClose) < int(cond.HistoryCloseCount) && tickTimeUnix.Before(lastTime) {
		return 0
	}
	var sellPrice float64
	var input quote.Quote
	input.Close = historyClose
	rsi, err := tickanalyze.GenerateRSI(input)
	if err != nil {
		logger.GetLogger().Errorf("GenerateRSI at GetSellPrice Stock: %s, Err: %s", tick.StockNum, err)
		return 0
	}
	switch {
	case tick.Close < stockutil.GetNewClose(originalOrderClose, -1) && rsi < cond.RsiLow:
		sellPrice = tick.Close
	case rsi > cond.RsiHigh:
		sellPrice = tick.Close
	case ManualSellMap.CheckStockExist(tick.StockNum):
		sellPrice = ManualSellMap.GetClose(tick.StockNum)
		if sellPrice == 0 {
			sellPrice = tick.Close
		}
	case tickTimeUnix.After(lastTime):
		sellPrice = tick.Close
	}
	return sellPrice
}

// CheckBuyOrderStatus CheckBuyOrderStatus
func CheckBuyOrderStatus(record traderecord.TradeRecord) {
	for {
		order, err := traderecord.GetOrderByOrderID(record.OrderID, db.GetAgent())
		if err != nil {
			logger.GetLogger().Error(err)
			continue
		}
		if order.Status == 4 {
			TradeQuota += record.BuyCost
			BuyOrderMap.DeleteByStockNum(record.StockNum)
			logger.GetLogger().WithFields(map[string]interface{}{
				"StockNum": order.StockNum,
				"Name":     order.StockName,
				"Quantity": order.Quantity,
				"Price":    order.Price,
			}).Info("Place Order Fail")
			return
		}
		if record.TradeTime.Add(30*time.Second).Before(time.Now()) && order.Status != 6 && order.Status != 5 {
			if err := Cancel(record.OrderID); err != nil {
				logger.GetLogger().Error(err)
				return
			}
			TradeQuota += record.BuyCost
			BuyOrderMap.DeleteByStockNum(record.StockNum)
			logger.GetLogger().WithFields(map[string]interface{}{
				"StockNum": order.StockNum,
				"Name":     order.StockName,
				"Quantity": order.Quantity,
				"Price":    order.Price,
			}).Info("Cancel Buy Order Success")
			return
		}
		if order.Status == 6 {
			FilledBuyOrderMap.Set(order)
			logger.GetLogger().WithFields(map[string]interface{}{
				"StockNum": order.StockNum,
				"Name":     order.StockName,
				"Quantity": order.Quantity,
				"Price":    order.Price,
			}).Info("Buy Stock Success")
			return
		}
		time.Sleep(1 * time.Second)
	}
}

// CheckSellOrderStatus CheckSellOrderStatus
func CheckSellOrderStatus(record traderecord.TradeRecord) {
	for {
		order, err := traderecord.GetOrderByOrderID(record.OrderID, db.GetAgent())
		if err != nil {
			logger.GetLogger().Error(err)
			continue
		}
		if order.Status == 4 {
			TradeQuota += record.BuyCost
			BuyOrderMap.DeleteByStockNum(record.StockNum)
			logger.GetLogger().WithFields(map[string]interface{}{
				"StockNum": order.StockNum,
				"Name":     order.StockName,
				"Quantity": order.Quantity,
				"Price":    order.Price,
			}).Info("Place Order Fail")
			return
		}
		if record.TradeTime.Add(45*time.Second).Before(time.Now()) && order.Status != 6 && order.Status != 5 {
			if err := Cancel(record.OrderID); err != nil {
				logger.GetLogger().Error(err)
				return
			}
			SellOrderMap.DeleteByStockNum(record.StockNum)
			logger.GetLogger().WithFields(map[string]interface{}{
				"StockNum": order.StockNum,
				"Name":     order.StockName,
				"Quantity": order.Quantity,
				"Price":    order.Price,
			}).Info("Cancel Sell Order Success")
			return
		}
		if order.Status == 6 {
			FilledSellOrderMap.Set(order)
			BuyOrderMap.DeleteByStockNum(record.StockNum)
			SellOrderMap.DeleteByStockNum(record.StockNum)
			if ManualSellMap.CheckStockExist(record.StockNum) {
				ManualSellMap.DeleteByStockNum(record.StockNum)
			}
			logger.GetLogger().WithFields(map[string]interface{}{
				"StockNum": order.StockNum,
				"Name":     order.StockName,
				"Quantity": order.Quantity,
				"Price":    order.Price,
			}).Info("Sell Stock Success")
			return
		}
		time.Sleep(1 * time.Second)
	}
}
