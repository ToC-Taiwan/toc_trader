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

// SellFirstOrderMap SellFirstOrderMap
var SellFirstOrderMap tradeRecordMutexMap

// BuyLaterOrderMap BuyLaterOrderMap
var BuyLaterOrderMap tradeRecordMutexMap

// FilledSellFirstOrderMap FilledSellFirstOrderMap
var FilledSellFirstOrderMap tradeRecordMutexMap

// FilledBuyLaterOrderMap FilledBuyLaterOrderMap
var FilledBuyLaterOrderMap tradeRecordMutexMap

// SellFirstBot SellFirstBot
func SellFirstBot(analyzeTick *analyzestreamtick.AnalyzeStreamTick) {
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
	}).Infof("StreamTick Analyze Sell First: %s %s %s", replaceDate, clockTime, analyzeTick.StockNum)

	buyCost := GetStockBuyCost(analyzeTick.Close, global.OneTimeQuantity)
	if global.TradeSwitch.SellFirst && SellFirstOrderMap.GetCount() < global.TradeSwitch.MeanTimeReverseTradeStockNum && TradeQuota-buyCost > 0 {
		if order, err := PlaceOrder(SellFirstAction, analyzeTick.StockNum, global.OneTimeQuantity, analyzeTick.Close); err != nil {
			logger.Logger.WithFields(map[string]interface{}{
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
		logger.Logger.Warn("Sell First Order is failed")
	}
}

// IsSellFirstPoint IsSellFirstPoint
func IsSellFirstPoint(analyzeTick *analyzestreamtick.AnalyzeStreamTick, cond simulationcond.AnalyzeCondition) bool {
	// closeChangeRatio := analyzeTick.CloseChangeRatio
	// if analyzeTick.OpenChangeRatio > cond.OpenChangeRatio || closeChangeRatio < cond.CloseChangeRatioLow || closeChangeRatio > cond.CloseChangeRatioHigh {
	// 	return false
	// }
	if analyzeTick.Volume < cond.Volume {
		return false
	}
	if analyzeTick.OutInRatio > cond.ReverseOutInRatio || analyzeTick.CloseDiff > cond.CloseDiff {
		return false
	}
	if analyzeTick.Rsi < float64(cond.ReverseRsiHigh) {
		return false
	}
	return true
}

// BuyLaterBot BuyLaterBot
func BuyLaterBot(ch chan *streamtick.StreamTick) {
	var historyClose []float64
	for {
		tick := <-ch
		historyClose = append(historyClose, tick.Close)
		if len(historyClose) > int(global.TickAnalyzeCondition.HistoryCloseCount) {
			historyClose = historyClose[1:]
		}
		filled, err := traderecord.CheckIsFilledByOrderID(SellFirstOrderMap.GetOrderID(tick.StockNum), global.GlobalDB)
		if err != nil {
			logger.Logger.Error(err)
			continue
		}
		if filled && !BuyLaterOrderMap.CheckStockExist(tick.StockNum) && global.TradeSwitch.Buy {
			originalOrderClose := SellFirstOrderMap.GetClose(tick.StockNum)
			buyPrice := GetBuyLaterPrice(tick, SellFirstOrderMap.GetTradeTime(tick.StockNum), historyClose, originalOrderClose, global.TickAnalyzeCondition)
			if buyPrice == 0 {
				continue
			} else if order, err := PlaceOrder(BuyAction, tick.StockNum, global.OneTimeQuantity, buyPrice); err != nil {
				logger.Logger.WithFields(map[string]interface{}{
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
				continue
			}
		}
	}
}

// GetBuyLaterPrice GetBuyLaterPrice
func GetBuyLaterPrice(tick *streamtick.StreamTick, tradeTime time.Time, historyClose []float64, originalOrderClose float64, cond simulationcond.AnalyzeCondition) float64 {
	tickTimeUnix := time.Unix(0, tick.TimeStamp)
	lastTime := time.Date(tickTimeUnix.Year(), tickTimeUnix.Month(), tickTimeUnix.Day(), 13, 0, 0, 0, time.Local)
	if len(historyClose) < int(cond.HistoryCloseCount) && tickTimeUnix.Before(lastTime) {
		return 0
	}
	var buyPrice float64
	var input quote.Quote
	input.Close = historyClose
	rsi, err := tickanalyze.GenerateRSI(input)
	if err != nil {
		logger.Logger.Errorf("GetBuyLaterPrice Stock: %s, Err: %s", tick.StockNum, err)
		return 0
	}
	switch {
	case tick.Close > stockutil.GetNewClose(originalOrderClose, 1) && rsi > float64(cond.ReverseRsiHigh):
		buyPrice = tick.Close
	case rsi < float64(cond.ReverseRsiLow):
		buyPrice = tick.Close
	case tickTimeUnix.After(lastTime):
		buyPrice = tick.Close
	default:
		buyPrice = 0
	}
	return buyPrice
}

// CheckSellFirstOrderStatus CheckSellFirstOrderStatus
func CheckSellFirstOrderStatus(record traderecord.TradeRecord) {
	for {
		time.Sleep(1 * time.Second)
		order, err := traderecord.GetOrderByOrderID(record.OrderID, global.GlobalDB)
		if err != nil {
			logger.Logger.Error(err)
			continue
		}
		if order.Status == 4 {
			TradeQuota += record.BuyCost
			SellFirstOrderMap.Delete(record.StockNum)
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
			SellFirstOrderMap.Delete(record.StockNum)
			logger.Logger.WithFields(map[string]interface{}{
				"StockNum": order.StockNum,
				"Name":     order.StockName,
				"Quantity": order.Quantity,
				"Price":    order.Price,
			}).Info("Cancel Sell First Order Success")
			return
		}
		if order.Status == 6 {
			FilledSellFirstOrderMap.Set(order)
			logger.Logger.WithFields(map[string]interface{}{
				"StockNum": order.StockNum,
				"Name":     order.StockName,
				"Quantity": order.Quantity,
				"Price":    order.Price,
			}).Info("Sell First Stock Success")
			return
		}
	}
}

// CheckBuyLaterOrderStatus CheckBuyLaterOrderStatus
func CheckBuyLaterOrderStatus(record traderecord.TradeRecord) {
	for {
		time.Sleep(1 * time.Second)
		order, err := traderecord.GetOrderByOrderID(record.OrderID, global.GlobalDB)
		if err != nil {
			logger.Logger.Error(err)
			continue
		}
		if order.Status == 4 {
			TradeQuota += record.BuyCost
			SellFirstOrderMap.Delete(record.StockNum)
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
			BuyLaterOrderMap.Delete(record.StockNum)
			logger.Logger.WithFields(map[string]interface{}{
				"StockNum": order.StockNum,
				"Name":     order.StockName,
				"Quantity": order.Quantity,
				"Price":    order.Price,
			}).Info("Cancel Buy Later Order Success")
			return
		}
		if order.Status == 6 {
			FilledBuyLaterOrderMap.Set(order)
			SellFirstOrderMap.Delete(record.StockNum)
			BuyLaterOrderMap.Delete(record.StockNum)
			logger.Logger.WithFields(map[string]interface{}{
				"StockNum": order.StockNum,
				"Name":     order.StockName,
				"Quantity": order.Quantity,
				"Price":    order.Price,
			}).Info("Buy Later Stock Success")
			return
		}
	}
}

func checkInSellFirstMap(stockNum string) bool {
	return FilledSellFirstOrderMap.CheckStockExist(stockNum) || SellFirstOrderMap.CheckStockExist(stockNum)
}