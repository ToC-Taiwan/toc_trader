// Package tradebot package tradebot
package tradebot

import (
	"time"

	"gitlab.tocraw.com/root/toc_trader/internal/db"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/traderecord"
)

// StatusFirstBack StatusFirstBack
var StatusFirstBack bool

// CheckOrderStatusLoop CheckOrderStatusLoop
func CheckOrderStatusLoop() {
	go showStatus()
	go tradeSwitch()
	tick := time.Tick(1*time.Second + 500*time.Millisecond)
	var initQuota bool
	for range tick {
		if err := FetchOrderStatus(); err != nil {
			logger.GetLogger().Error(err)
			continue
		}
		if !initQuota && StatusFirstBack {
			if err := InitStartUpQuota(); err != nil {
				panic(err)
			}
			logger.GetLogger().Warnf("Initial Quota: %d", TradeQuota)
			dbOrder, err := traderecord.GetAllorderByDayTime(global.TradeDay, db.GetAgent())
			if err != nil {
				logger.GetLogger().Error(err)
				continue
			}
			initBalance(dbOrder)
			initQuota = true
		}
	}
}

// showStatus showStatus
func showStatus() {
	tick := time.Tick(60 * time.Second)
	for range tick {
		if isCurrentOrderAllFinished() && time.Now().Before(global.TradeDayOutEndTime.Add(2*time.Hour)) {
			balance := FilledSellOrderMap.GetTotalSellCost() + FilledSellFirstOrderMap.GetTotalSellCost() - FilledBuyLaterOrderMap.GetTotalBuyCost() - FilledBuyOrderMap.GetTotalBuyCost()
			back := FilledBuyOrderMap.GetTotalCostBack() + FilledSellOrderMap.GetTotalCostBack() + FilledSellFirstOrderMap.GetTotalCostBack() + FilledBuyLaterOrderMap.GetTotalCostBack()
			logger.GetLogger().WithFields(map[string]interface{}{
				"Current":         BuyOrderMap.GetCount(),
				"Maximum":         global.TradeSwitch.MeanTimeTradeStockNum,
				"TradeQuota":      TradeQuota,
				"OriginalBalance": balance,
				"Back":            back,
				"Real":            balance + back,
			}).Info("Current Trade Status")
			if balance < -1000 && (global.TradeSwitch.Buy || global.TradeSwitch.SellFirst) {
				global.TradeSwitch.Buy = false
				global.TradeSwitch.SellFirst = false
				logger.GetLogger().Warn("Enable buy and Sell first are all OFF because too...")
			}
		}
	}
}

func tradeSwitch() {
	tick := time.Tick(20 * time.Second)
	for range tick {
		if time.Now().After(global.TradeDayInEndTime) && (global.TradeSwitch.Buy || global.TradeSwitch.SellFirst) {
			global.TradeSwitch.Buy = false
			global.TradeSwitch.SellFirst = false
			logger.GetLogger().Warn("Enable buy and Sell first are all OFF")
		}
	}
}

func isCurrentOrderAllFinished() bool {
	if BuyOrderMap.GetCount() != 0 || SellFirstOrderMap.GetCount() != 0 {
		return false
	}
	if FilledSellOrderMap.GetCount() == 0 && FilledBuyLaterOrderMap.GetCount() == 0 {
		return false
	}
	return FilledBuyOrderMap.GetCount() == FilledSellOrderMap.GetCount() && FilledSellFirstOrderMap.GetCount() == FilledBuyLaterOrderMap.GetCount()
}

func initBalance(orders []traderecord.TradeRecord) {
	var tmp []string
	for _, val := range orders {
		record := traderecord.TradeRecord{
			StockNum:  val.StockNum,
			StockName: global.AllStockNameMap.GetName(val.StockNum),
			Action:    val.Action,
			Price:     val.Price,
			Quantity:  val.Quantity,
			Status:    val.Status,
			OrderID:   val.OrderID,
			TradeTime: time.Now(),
		}
		if val.Action == 1 && val.Status == 6 && !FilledBuyOrderMap.CheckStockExist(val.StockNum) {
			FilledBuyOrderMap.Set(record)
			tmp = append(tmp, val.StockNum)
		} else if val.Action == 2 && val.Status == 6 && !FilledSellOrderMap.CheckStockExist(val.StockNum) {
			FilledSellOrderMap.Set(record)
		}
	}
	for _, v := range tmp {
		buyOrder := FilledBuyOrderMap.GetRecordByStockNum(v)
		logger.GetLogger().WithFields(map[string]interface{}{
			"StockNum": buyOrder.StockNum,
			"Name":     buyOrder.StockName,
			"Quantity": buyOrder.Quantity,
			"Price":    buyOrder.Price,
		}).Warn("Filled Buy Order")
		sellOrder := FilledSellOrderMap.GetRecordByStockNum(v)
		logger.GetLogger().WithFields(map[string]interface{}{
			"StockNum": sellOrder.StockNum,
			"Name":     sellOrder.StockName,
			"Quantity": sellOrder.Quantity,
			"Price":    sellOrder.Price,
		}).Warn("Filled Sell Order")
	}
}
