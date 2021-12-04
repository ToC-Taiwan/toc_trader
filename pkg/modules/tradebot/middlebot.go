// Package tradebot package tradebot
package tradebot

import (
	"time"

	"gitlab.tocraw.com/root/toc_trader/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/analyzestreamtick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/biasrate"
)

var (
	tradeInWaitTime  time.Duration = 15 * time.Second
	tradeOutWaitTime time.Duration = 30 * time.Second
)

// BuyAgent BuyAgent
func BuyAgent(ch chan *analyzestreamtick.AnalyzeStreamTick) {
	for {
		analyzeTick := <-ch
		if BuyOrderMap.CheckStockExist(analyzeTick.StockNum) || FilledBuyOrderMap.CheckStockExist(analyzeTick.StockNum) || SellFirstOrderMap.CheckStockExist(analyzeTick.StockNum) {
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
		if SellFirstOrderMap.CheckStockExist(analyzeTick.StockNum) || FilledSellFirstOrderMap.CheckStockExist(analyzeTick.StockNum) || BuyOrderMap.CheckStockExist(analyzeTick.StockNum) {
			continue
		}
		if IsSellFirstPoint(analyzeTick, global.ReverseCond) {
			if global.TradeSwitch.SellFirst {
				SellFirstBot(analyzeTick)
			}
		}
	}
}

// GetQuantityByTradeDay GetQuantityByTradeDay
func GetQuantityByTradeDay(stockNum, tradeDay string, tradeType global.TradeType) int64 {
	var quantity int64 = 1
	biasRate := biasrate.StockBiasRateMap.GetBiasRate(stockNum, tradeDay)
	if tradeType == global.ForwardTrade && biasRate > 4 {
		quantity = 2
	}
	if tradeType == global.ReverseTrade && biasRate < -4 {
		quantity = 2
	}
	return quantity
}
