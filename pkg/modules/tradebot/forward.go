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

// BuyBot BuyBot
func BuyBot(ch chan *analyzestreamtick.AnalyzeStreamTick, cond simulationcond.AnalyzeCondition) {
	for {
		analyzeTick := <-ch
		if checkInBuyMap(analyzeTick.StockNum) || checkInSellFirstMap(analyzeTick.StockNum) {
			continue
		}
		if !IsBuyPoint(analyzeTick, cond) {
			if IsSellFirstPoint(analyzeTick, cond) {
				go SellFirstBot(analyzeTick)
			}
			continue
		}

		name := global.AllStockNameMap.GetName(analyzeTick.StockNum)
		outSum := analyzeTick.OutSum
		inSum := analyzeTick.InSum
		closeChangeRatio := analyzeTick.CloseChangeRatio
		outInRatio := analyzeTick.OutInRatio

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

		buyCost := GetStockBuyCost(analyzeTick.Close, global.OneTimeQuantity)
		if global.TradeSwitch.Buy && BuyOrderMap.GetCount() < global.TradeSwitch.MeanTimeTradeStockNum && TradeQuota-buyCost > 0 {
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
				continue
			}
			logger.Logger.Warn("Buy Order is failed")
		}
	}
}

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

// SellBot SellBot
func SellBot(ch chan *streamtick.StreamTick, cond simulationcond.AnalyzeCondition) {
	var historyClose []float64
	for {
		tick := <-ch
		historyClose = append(historyClose, tick.Close)
		if len(historyClose) > int(cond.HistoryCloseCount) {
			historyClose = historyClose[1:]
		}
		filled, err := traderecord.CheckIsFilledByOrderID(BuyOrderMap.GetOrderID(tick.StockNum), global.GlobalDB)
		if err != nil {
			logger.Logger.Error(err)
			continue
		}
		if filled && !SellOrderMap.CheckStockExist(tick.StockNum) && global.TradeSwitch.Sell {
			originalOrderClose := BuyOrderMap.GetClose(tick.StockNum)
			sellPrice := GetSellPrice(tick, BuyOrderMap.GetTradeTime(tick.StockNum), historyClose, originalOrderClose, cond)
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
func GetSellPrice(tick *streamtick.StreamTick, tradeTime time.Time, historyClose []float64, originalOrderClose float64, cond simulationcond.AnalyzeCondition) float64 {
	tickTimeUnix := time.Unix(0, tick.TimeStamp)
	lastTime := time.Date(tickTimeUnix.Year(), tickTimeUnix.Month(), tickTimeUnix.Day(), global.TradeEndHour, global.TradeEndMinute, 0, 0, time.Local)
	if len(historyClose) < int(cond.HistoryCloseCount) && tickTimeUnix.Before(lastTime) {
		return 0
	}
	var sellPrice float64
	var input quote.Quote
	input.Close = historyClose
	rsi, err := tickanalyze.GenerateRSI(input)
	if err != nil {
		logger.Logger.Errorf("GetSellPrice Stock: %s, Err: %s", tick.StockNum, err)
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
		if record.TradeTime.Add(30*time.Second).Before(time.Now()) && order.Status != 6 && order.Status != 5 {
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
		if record.TradeTime.Add(30*time.Second).Before(time.Now()) && order.Status != 6 && order.Status != 5 {
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

func checkInBuyMap(stockNum string) bool {
	return FilledBuyOrderMap.CheckStockExist(stockNum) || BuyOrderMap.CheckStockExist(stockNum)
}
