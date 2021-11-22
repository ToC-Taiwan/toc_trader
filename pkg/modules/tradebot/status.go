// Package tradebot package tradebot
package tradebot

import (
	"time"

	"gitlab.tocraw.com/root/toc_trader/internal/database"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/balance"
)

// CheckOrderStatusLoop CheckOrderStatusLoop
func CheckOrderStatusLoop() {
	go showStatus()
	go tradeSwitch()

	for range time.Tick(1*time.Second + 500*time.Millisecond) {
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
			forward := FilledSellOrderMap.GetTotalSellCost() - FilledBuyOrderMap.GetTotalBuyCost()
			reverse := FilledSellFirstOrderMap.GetTotalSellCost() - FilledBuyLaterOrderMap.GetTotalBuyCost()
			totalCount := FilledSellOrderMap.GetCount() + FilledBuyLaterOrderMap.GetCount()
			discount := FilledBuyOrderMap.GetTotalCostBack() + FilledSellOrderMap.GetTotalCostBack() + FilledSellFirstOrderMap.GetTotalCostBack() + FilledBuyLaterOrderMap.GetTotalCostBack()
			sum := balance.Balance{
				TradeDay:        global.TradeDay,
				TradeCount:      int64(totalCount),
				Forward:         forward,
				Reverse:         reverse,
				OriginalBalance: forward + reverse,
				Discount:        discount,
				Total:           forward + reverse + discount,
			}
			if err := balance.InsertOrUpdate(sum, database.GetAgent()); err != nil {
				logger.GetLogger().Error(err)
				continue
			}
			if tmpBalance != sum.Total {
				tmpBalance = sum.Total
				logger.GetLogger().WithFields(map[string]interface{}{
					"Current":         BuyOrderMap.GetCount(),
					"Maximum":         global.TradeSwitch.MeanTimeTradeStockNum,
					"TradeQuota":      TradeQuota,
					"OriginalBalance": sum.OriginalBalance,
					"Back":            sum.Discount,
					"Real":            sum.Total,
				}).Info("TradeStatus")
			}

			if sum.Total < -500 && (global.TradeSwitch.Buy || global.TradeSwitch.SellFirst) {
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
