// Package tradebot package tradebot
package tradebot

import (
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/analyzestreamtick"
)

// BuyAgent BuyAgent
func BuyAgent(ch chan *analyzestreamtick.AnalyzeStreamTick) {
	for {
		analyzeTick := <-ch
		if checkInBuyMap(analyzeTick.StockNum) {
			continue
		}
		if IsBuyPoint(analyzeTick, global.ForwardCond) && global.TradeSwitch.Buy {
			go BuyBot(analyzeTick)
		}
	}
}

// SellFirstAgent SellFirstAgent
func SellFirstAgent(ch chan *analyzestreamtick.AnalyzeStreamTick) {
	for {
		analyzeTick := <-ch
		if checkInSellFirstMap(analyzeTick.StockNum) {
			continue
		}
		if IsSellFirstPoint(analyzeTick, global.ReverseCond) && global.TradeSwitch.SellFirst {
			go SellFirstBot(analyzeTick)
		}
	}
}

func checkInBuyMap(stockNum string) bool {
	return FilledBuyOrderMap.CheckStockExist(stockNum) || BuyOrderMap.CheckStockExist(stockNum)
}

func checkInSellFirstMap(stockNum string) bool {
	return FilledSellFirstOrderMap.CheckStockExist(stockNum) || SellFirstOrderMap.CheckStockExist(stockNum)
}
