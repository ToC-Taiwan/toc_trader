// Package tradebot package tradebot
package tradebot

import (
	"time"

	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
)

// CheckOrderStatusLoop CheckOrderStatusLoop
func CheckOrderStatusLoop() {
	go showStatus()
	go tradeSwitch()
	tick := time.Tick(1*time.Second + 500*time.Millisecond)
	for range tick {
		if err := FetchOrderStatus(); err != nil {
			logger.GetLogger().Error(err)
		}
	}
}

// showStatus showStatus
func showStatus() {
	tick := time.Tick(60 * time.Second)
	var tmpBalance int64
	for range tick {
		if isCurrentOrderAllFinished() {
			balance := FilledSellOrderMap.GetTotalSellCost() + FilledSellFirstOrderMap.GetTotalSellCost() - FilledBuyLaterOrderMap.GetTotalBuyCost() - FilledBuyOrderMap.GetTotalBuyCost()
			back := FilledBuyOrderMap.GetTotalCostBack() + FilledSellOrderMap.GetTotalCostBack() + FilledSellFirstOrderMap.GetTotalCostBack() + FilledBuyLaterOrderMap.GetTotalCostBack()
			if tmpBalance != (balance + back) {
				tmpBalance = (balance + back)
				logger.GetLogger().WithFields(map[string]interface{}{
					"Current":         BuyOrderMap.GetCount(),
					"Maximum":         global.TradeSwitch.MeanTimeTradeStockNum,
					"TradeQuota":      TradeQuota,
					"OriginalBalance": balance,
					"Back":            back,
					"Real":            balance + back,
				}).Info("TradeStatus")
			}
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

// CheckIsOpenTime CheckIsOpenTime
func CheckIsOpenTime() bool {
	starTime := global.TradeDay.Add(1 * time.Hour)
	if time.Now().After(starTime) && time.Now().Before(global.TradeDayInEndTime) {
		return true
	}
	return false
}
