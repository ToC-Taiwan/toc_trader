// Package tradebot package tradebot
package tradebot

import (
	"time"

	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/analyzestreamtick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/biasrate"
)

var (
	tradeInWaitTime  time.Duration = 15 * time.Second
	tradeOutWaitTime time.Duration = 45 * time.Second
)

// BuyAgent BuyAgent
func BuyAgent(ch chan *analyzestreamtick.AnalyzeStreamTick) {
	for {
		analyzeTick := <-ch
		if checkAlreadyTraded(analyzeTick.StockNum) {
			continue
		}
		if IsBuyPoint(analyzeTick, global.ForwardCond) {
			if global.TradeSwitch.Buy {
				BuyBot(analyzeTick)
			}
		}
	}
}

// SellFirstAgent SellFirstAgent
func SellFirstAgent(ch chan *analyzestreamtick.AnalyzeStreamTick) {
	for {
		analyzeTick := <-ch
		if checkAlreadyTraded(analyzeTick.StockNum) {
			continue
		}
		if IsSellFirstPoint(analyzeTick, global.ReverseCond) {
			if global.TradeSwitch.SellFirst {
				SellFirstBot(analyzeTick)
			}
		}
	}
}

func checkAlreadyTraded(stockNum string) bool {
	if BuyOrderMap.CheckStockExist(stockNum) || FilledBuyOrderMap.CheckStockExist(stockNum) {
		return true
	}
	if SellFirstOrderMap.CheckStockExist(stockNum) || FilledSellFirstOrderMap.CheckStockExist(stockNum) {
		return true
	}
	return false
}

// GetQuantityByTradeDay GetQuantityByTradeDay
func GetQuantityByTradeDay(stockNum, tradeDay string) int64 {
	var quantity int64 = 2
	biasRate := biasrate.StockBiasRateMap.GetBiasRate(stockNum, tradeDay)
	if biasRate < 10 && biasRate > -10 {
		quantity = 1
	}
	// if biasRate == 0 {
	// 	return 0
	// }
	return quantity
}
