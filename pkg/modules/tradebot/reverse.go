// Package tradebot package tradebot
package tradebot

import (
	"time"

	"github.com/markcheno/go-quote"
	"gitlab.tocraw.com/root/toc_trader/external/sinopacsrv"
	"gitlab.tocraw.com/root/toc_trader/internal/db"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/internal/stockutil"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/analyzestreamtick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/streamtick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/traderecord"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tickanalyze"
)

// SellFirstOrderMap SellFirstOrderMap
var SellFirstOrderMap tradeRecordMutexMap

// BuyLaterOrderMap BuyLaterOrderMap
var BuyLaterOrderMap tradeRecordMutexMap

// FilledSellFirstOrderMap FilledSellFirstOrderMap
var FilledSellFirstOrderMap tradeRecordMutexMap

// FilledBuyLaterOrderMap FilledBuyLaterOrderMap
var FilledBuyLaterOrderMap tradeRecordMutexMap

// ManualBuyLaterMap ManualBuyLaterMap
var ManualBuyLaterMap tradeRecordMutexMap

// IsSellFirstPoint IsSellFirstPoint
func IsSellFirstPoint(analyzeTick *analyzestreamtick.AnalyzeStreamTick, cond simulationcond.AnalyzeCondition) bool {
	// closeChangeRatio := analyzeTick.CloseChangeRatio
	// if analyzeTick.OpenChangeRatio > cond.OpenChangeRatio || closeChangeRatio < cond.CloseChangeRatioLow || closeChangeRatio > cond.CloseChangeRatioHigh {
	// 	return false
	// }
	if analyzeTick.Volume < cond.VolumePerSecond*int64(analyzeTick.TotalTime) {
		return false
	}
	if analyzeTick.OutInRatio > cond.ReverseOutInRatio || analyzeTick.CloseDiff > cond.CloseDiff {
		return false
	}
	if analyzeTick.Rsi < cond.ReverseRsiHigh {
		return false
	}
	return true
}

// SellFirstBot SellFirstBot
func SellFirstBot(analyzeTick *analyzestreamtick.AnalyzeStreamTick) {
	buyCost := GetStockBuyCost(analyzeTick.Close, global.OneTimeQuantity)
	if SellFirstOrderMap.GetCount() < global.TradeSwitch.MeanTimeReverseTradeStockNum && TradeQuota-buyCost > 0 {
		if order, err := PlaceOrder(SellFirstAction, analyzeTick.StockNum, global.OneTimeQuantity, analyzeTick.Close); err != nil {
			logger.GetLogger().WithFields(map[string]interface{}{
				"Msg":      err,
				"StockNum": analyzeTick.StockNum,
				"Quantity": global.OneTimeQuantity,
				"Price":    analyzeTick.Close,
			}).Error("Sell First fail")
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
			SellFirstOrderMap.Set(record)
			go CheckSellFirstOrderStatus(record)
		}
	} else {
		logger.GetLogger().Warn("Over MeanTimeReverseTradeStockNum or Quota")
	}
}

// BuyLaterBot BuyLaterBot
func BuyLaterBot(ch chan *streamtick.StreamTick, cond simulationcond.AnalyzeCondition, historyClosePrt *[]float64) {
	// var historyClose []float64
	var filled bool
	for {
		tick := <-ch
		historyClose := *historyClosePrt
		if !filled {
			if tmpFilled, err := traderecord.CheckIsFilledByOrderID(SellFirstOrderMap.GetOrderIDByStockNum(tick.StockNum), db.GetAgent()); err != nil {
				logger.GetLogger().Error(err)
				continue
			} else if tmpFilled {
				filled = true
			}
		} else if !BuyLaterOrderMap.CheckStockExist(tick.StockNum) {
			originalOrderClose := SellFirstOrderMap.GetClose(tick.StockNum)
			buyPrice := GetBuyLaterPrice(tick, SellFirstOrderMap.GetTradeTime(tick.StockNum), historyClose, originalOrderClose, cond)
			if buyPrice == 0 {
				continue
			} else if order, err := PlaceOrder(BuyAction, tick.StockNum, global.OneTimeQuantity, buyPrice); err != nil {
				logger.GetLogger().WithFields(map[string]interface{}{
					"Msg":      err,
					"Stock":    tick.StockNum,
					"Quantity": global.OneTimeQuantity,
					"Price":    buyPrice,
				}).Error("Buy Later fail")
			} else if order.OrderID != "" && order.Status != traderecord.Failed {
				record := traderecord.TradeRecord{
					StockNum:  tick.StockNum,
					Price:     buyPrice,
					Quantity:  global.OneTimeQuantity,
					Action:    int64(SellAction),
					TradeTime: time.Unix(0, tick.TimeStamp),
					OrderID:   order.OrderID,
				}
				BuyLaterOrderMap.Set(record)
				go CheckBuyLaterOrderStatus(record)
			}
		}
	}
}

// GetBuyLaterPrice GetBuyLaterPrice
func GetBuyLaterPrice(tick *streamtick.StreamTick, tradeTime time.Time, historyClose []float64, originalOrderClose float64, cond simulationcond.AnalyzeCondition) float64 {
	if tick.PctChg < -9.9 {
		return tick.Close
	}
	tickTimeUnix := time.Unix(0, tick.TimeStamp)
	lastTime := time.Date(tickTimeUnix.Year(), tickTimeUnix.Month(), tickTimeUnix.Day(), global.TradeOutEndHour, global.TradeOutEndMinute, 0, 0, time.Local)
	if len(historyClose) < int(cond.HistoryCloseCount) && tickTimeUnix.Before(lastTime) {
		return 0
	}
	var buyPrice float64
	var input quote.Quote
	input.Close = historyClose
	rsi, err := tickanalyze.GenerateRSI(input)
	if err != nil {
		logger.GetLogger().Errorf("GenerateRSI at GetBuyLaterPrice Stock: %s, Err: %s", tick.StockNum, err)
		return 0
	}
	switch {
	case tick.Close > stockutil.GetNewClose(originalOrderClose, 1) && rsi > cond.ReverseRsiHigh:
		buyPrice = tick.Close
	case rsi < cond.ReverseRsiLow:
		buyPrice = tick.Close
	case ManualBuyLaterMap.CheckStockExist(tick.StockNum):
		buyPrice = ManualBuyLaterMap.GetClose(tick.StockNum)
		if buyPrice == 0 {
			buyPrice = tick.Close
		}
	case tickTimeUnix.After(lastTime):
		buyPrice = tick.Close
	}
	return buyPrice
}

// CheckSellFirstOrderStatus CheckSellFirstOrderStatus
func CheckSellFirstOrderStatus(record traderecord.TradeRecord) {
	var cancelAlready bool
	for {
		time.Sleep(time.Second)
		order, err := traderecord.GetOrderByOrderID(record.OrderID, db.GetAgent())
		if err != nil {
			logger.GetLogger().Error(err)
			continue
		}
		if order.Status == 4 || order.Status == 5 {
			TradeQuota += record.BuyCost
			SellFirstOrderMap.DeleteByStockNum(record.StockNum)
			logger.GetLogger().WithFields(map[string]interface{}{
				"StockNum": order.StockNum, "Name": order.StockName, "Quantity": order.Quantity, "Price": order.Price, "Status": order.Status,
			}).Info("CheckSellFirstOrderStatus: Order Fail or Canceled")
			return
		}
		if order.Status == 6 {
			FilledSellFirstOrderMap.Set(order)
			logger.GetLogger().WithFields(map[string]interface{}{
				"StockNum": order.StockNum, "Name": order.StockName, "Quantity": order.Quantity, "Price": order.Price,
			}).Info("Sell First Stock Success")
			return
		}
		if record.TradeTime.Add(tradeInWaitTime).Before(time.Now()) && !cancelAlready {
			if err := Cancel(record.OrderID); err != nil {
				if err.Error() == sinopacsrv.StatusAlready {
					cancelAlready = true
					continue
				}
				logger.GetLogger().Error(err)
				continue
			}
		}
	}
}

// CheckBuyLaterOrderStatus CheckBuyLaterOrderStatus
func CheckBuyLaterOrderStatus(record traderecord.TradeRecord) {
	var cancelAlready bool
	for {
		time.Sleep(time.Second)
		order, err := traderecord.GetOrderByOrderID(record.OrderID, db.GetAgent())
		if err != nil {
			logger.GetLogger().Error(err)
			continue
		}
		if order.Status == 4 || order.Status == 5 {
			TradeQuota += record.BuyCost
			BuyLaterOrderMap.DeleteByStockNum(record.StockNum)
			logger.GetLogger().WithFields(map[string]interface{}{
				"StockNum": order.StockNum, "Name": order.StockName, "Quantity": order.Quantity, "Price": order.Price, "Status": order.Status,
			}).Info("CheckBuyLaterOrderStatus: Order Fail or Canceled")
			return
		}
		if order.Status == 6 {
			FilledBuyLaterOrderMap.Set(order)
			SellFirstOrderMap.DeleteByStockNum(record.StockNum)
			BuyLaterOrderMap.DeleteByStockNum(record.StockNum)
			if ManualBuyLaterMap.CheckStockExist(record.StockNum) {
				ManualBuyLaterMap.DeleteByStockNum(record.StockNum)
			}
			logger.GetLogger().WithFields(map[string]interface{}{
				"StockNum": order.StockNum, "Name": order.StockName, "Quantity": order.Quantity, "Price": order.Price,
			}).Info("Buy Later Stock Success")
			return
		}
		if record.TradeTime.Add(tradeOutWaitTime).Before(time.Now()) && !cancelAlready {
			if err := Cancel(record.OrderID); err != nil {
				if err.Error() == sinopacsrv.StatusAlready {
					cancelAlready = true
					continue
				}
				logger.GetLogger().Error(err)
				continue
			}
		}
	}
}
