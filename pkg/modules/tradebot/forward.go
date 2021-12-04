// Package tradebot package tradebot
package tradebot

import (
	"time"

	"gitlab.tocraw.com/root/toc_trader/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/database"
	"gitlab.tocraw.com/root/toc_trader/pkg/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/analyzestreamtick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/streamtick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/traderecord"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tickanalyze"
	"gitlab.tocraw.com/root/toc_trader/pkg/sinopacapi"
	"gitlab.tocraw.com/root/toc_trader/pkg/utils"
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
func BuyBot(analyzeTick *analyzestreamtick.AnalyzeStreamTick) {
	quantity := GetQuantityByTradeDay(analyzeTick.StockNum, global.TradeDay.Format(global.ShortTimeLayout), global.ForwardTrade)
	if quantity == 0 {
		logger.GetLogger().Warnf("%s quantity is 0", analyzeTick.StockNum)
		return
	}
	buyCost := GetStockBuyCost(analyzeTick.Close, quantity)
	if BuyOrderMap.GetCount() < global.TradeSwitch.MeanTimeTradeStockNum && TradeQuota-buyCost > 0 {
		newOrder := sinopacapi.Order{
			StockNum:  analyzeTick.StockNum,
			Price:     analyzeTick.Close,
			Quantity:  quantity,
			OrderID:   "",
			Action:    sinopacapi.ActionBuy,
			TradeTime: time.Unix(0, analyzeTick.TimeStamp),
		}
		if orderRes, err := sinopacapi.GetAgent().PlaceOrder(newOrder); err != nil {
			logger.GetLogger().WithFields(map[string]interface{}{
				"Msg":      err,
				"StockNum": analyzeTick.StockNum,
				"Quantity": quantity,
				"Price":    analyzeTick.Close,
			}).Error("Buy fail")
		} else if orderRes.OrderID != "" && orderRes.Status != traderecord.Failed {
			TradeQuota -= buyCost
			record := traderecord.TradeRecord{
				StockNum:  analyzeTick.StockNum,
				Price:     analyzeTick.Close,
				Quantity:  quantity,
				Action:    int64(sinopacapi.ActionBuy),
				BuyCost:   buyCost,
				TradeTime: time.Unix(0, analyzeTick.TimeStamp),
				OrderID:   orderRes.OrderID,
			}
			BuyOrderMap.Set(record)
			go CheckBuyOrderStatus(record)
		}
	} else {
		logger.GetLogger().Warn("Over MeanTimeTradeStockNum or Quota")
	}
}

// SellBot SellBot
func SellBot(ch chan *streamtick.StreamTick, cond simulationcond.AnalyzeCondition, historyClosePtr *[]float64) {
	var maxClose float64
	for {
		tick := <-ch
		if maxClose == 0 {
			maxClose = tick.Close
		} else if tick.Close > maxClose {
			maxClose = tick.Close
		}
		if !SellOrderMap.CheckStockExist(tick.StockNum) {
			originalOrderClose := BuyOrderMap.GetClose(tick.StockNum)
			quantity := GetQuantityByTradeDay(tick.StockNum, global.TradeDay.Format(global.ShortTimeLayout), global.ForwardTrade)
			sellPrice := GetSellPrice(tick, BuyOrderMap.GetTradeTime(tick.StockNum), *historyClosePtr, originalOrderClose, maxClose, cond)
			if sellPrice == 0 {
				continue
			}
			newOrder := sinopacapi.Order{
				StockNum:  tick.StockNum,
				Price:     sellPrice,
				Quantity:  quantity,
				OrderID:   "",
				Action:    sinopacapi.ActionSell,
				TradeTime: time.Unix(0, tick.TimeStamp),
			}
			if orderRes, err := sinopacapi.GetAgent().PlaceOrder(newOrder); err != nil {
				logger.GetLogger().WithFields(map[string]interface{}{
					"Msg":      err,
					"Stock":    tick.StockNum,
					"Quantity": quantity,
					"Price":    sellPrice,
				}).Error("Sell fail")
			} else if orderRes.OrderID != "" && orderRes.Status != traderecord.Failed {
				record := traderecord.TradeRecord{
					StockNum:  tick.StockNum,
					Price:     sellPrice,
					Quantity:  quantity,
					Action:    int64(sinopacapi.ActionSell),
					TradeTime: time.Unix(0, tick.TimeStamp),
					OrderID:   orderRes.OrderID,
				}
				SellOrderMap.Set(record)
				go CheckSellOrderStatus(record)
			}
		}
	}
}

// IsBuyPoint IsBuyPoint
func IsBuyPoint(analyzeTick *analyzestreamtick.AnalyzeStreamTick, cond simulationcond.AnalyzeCondition) bool {
	closeChangeRatio := analyzeTick.CloseChangeRatio
	if analyzeTick.Volume < cond.VolumePerSecond*int64(analyzeTick.TotalTime) {
		return false
	}
	if analyzeTick.OpenChangeRatio > cond.OpenChangeRatio || closeChangeRatio > cond.CloseChangeRatioHigh || closeChangeRatio < cond.CloseChangeRatioLow {
		return false
	}
	if analyzeTick.OutInRatio < cond.ForwardOutInRatio {
		return false
	}
	return true
}

// GetSellPrice GetSellPrice
func GetSellPrice(tick *streamtick.StreamTick, tradeTime time.Time, historyClose []float64, originalOrderClose, maxClose float64, cond simulationcond.AnalyzeCondition) float64 {
	if tick.Close >= utils.GetMaxByOpen(tick.Open) {
		return tick.Close
	}
	tickTimeUnix := time.Unix(0, tick.TimeStamp)
	lastTime := time.Date(tickTimeUnix.Year(), tickTimeUnix.Month(), tickTimeUnix.Day(), global.TradeOutEndHour, global.TradeOutEndMinute, 0, 0, time.Local)
	if len(historyClose) < int(cond.HistoryCloseCount) && tickTimeUnix.Before(lastTime) {
		return 0
	}
	var sellPrice float64
	rsiHighStatus := tickanalyze.GetForwardRSIStatus(historyClose, cond.RsiHigh)
	switch {
	case tick.Close/originalOrderClose < 0.99:
		sellPrice = tick.Close
	case rsiHighStatus && tick.Close > originalOrderClose:
		sellPrice = tick.Close
	case tickTimeUnix.After(lastTime):
		sellPrice = tick.Close
	case ManualSellMap.CheckStockExist(tick.StockNum):
		sellPrice = ManualSellMap.GetClose(tick.StockNum)
		if sellPrice == 0 {
			sellPrice = tick.Close
		}
	}
	holdTime := 45 * int64(time.Minute)
	if sellPrice == 0 && tradeTime.Add(time.Duration(holdTime)).Before(tickTimeUnix) {
		if tick.Close < utils.GetNewClose(maxClose, -2) && tick.Close > originalOrderClose {
			sellPrice = tick.Close
		}
	}
	return sellPrice
}

// CheckBuyOrderStatus CheckBuyOrderStatus
func CheckBuyOrderStatus(record traderecord.TradeRecord) {
	for {
		time.Sleep(time.Second)
		order, err := traderecord.GetOrderByOrderID(record.OrderID, database.GetAgent())
		if err != nil {
			logger.GetLogger().Error(err)
			continue
		}
		if order.Status == 5 {
			TradeQuota += record.BuyCost
			BuyOrderMap.DeleteByStockNum(record.StockNum)
			logger.GetLogger().WithFields(map[string]interface{}{
				"StockNum": order.StockNum, "Name": order.StockName, "Quantity": order.Quantity, "Price": order.Price, "Status": order.Status,
			}).Info("CheckBuyOrderStatus: Order Fail or Canceled")
			return
		}
		if order.Status == 6 {
			FilledBuyOrderMap.Set(order)
			logger.GetLogger().WithFields(map[string]interface{}{
				"StockNum": order.StockNum, "Name": order.StockName, "Quantity": order.Quantity, "Price": order.Price,
			}).Info("Buy Stock Success")
			return
		}
		if record.TradeTime.Add(tradeInWaitTime).Before(time.Now()) {
			if err := sinopacapi.GetAgent().CancelOrder(record.OrderID); err != nil {
				logger.GetLogger().WithFields(map[string]interface{}{
					"StockNum": order.StockNum,
					"Name":     order.StockName,
					"Quantity": order.Quantity,
					"Price":    order.Price,
					"Error":    err,
				}).Error("Cancel Fail")
			}
			time.Sleep(5 * time.Second)
		}
	}
}

// CheckSellOrderStatus CheckSellOrderStatus
func CheckSellOrderStatus(record traderecord.TradeRecord) {
	for {
		time.Sleep(time.Second)
		order, err := traderecord.GetOrderByOrderID(record.OrderID, database.GetAgent())
		if err != nil {
			logger.GetLogger().Error(err)
			continue
		}
		if order.Status == 5 {
			TradeQuota += record.BuyCost
			SellOrderMap.DeleteByStockNum(record.StockNum)
			logger.GetLogger().WithFields(map[string]interface{}{
				"StockNum": order.StockNum, "Name": order.StockName, "Quantity": order.Quantity, "Price": order.Price, "Status": order.Status,
			}).Info("CheckSellOrderStatus: Order Fail or Canceled")
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
				"StockNum": order.StockNum, "Name": order.StockName, "Quantity": order.Quantity, "Price": order.Price,
			}).Info("Sell Stock Success")
			return
		}
		if record.TradeTime.Add(tradeOutWaitTime).Before(time.Now()) {
			if err := sinopacapi.GetAgent().CancelOrder(record.OrderID); err != nil {
				logger.GetLogger().WithFields(map[string]interface{}{
					"StockNum": order.StockNum,
					"Name":     order.StockName,
					"Quantity": order.Quantity,
					"Price":    order.Price,
					"Error":    err,
				}).Error("Cancel Fail")
			}
			time.Sleep(5 * time.Second)
		}
	}
}
