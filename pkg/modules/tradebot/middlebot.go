// Package tradebot package tradebot
package tradebot

import (
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/analyzestreamtick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
)

// TradeAgent TradeAgent
func TradeAgent(ch chan *analyzestreamtick.AnalyzeStreamTick, cond simulationcond.AnalyzeCondition) {
	for {
		analyzeTick := <-ch
		if checkInBuyMap(analyzeTick.StockNum) || checkInSellFirstMap(analyzeTick.StockNum) {
			continue
		}
		if IsBuyPoint(analyzeTick, cond) && global.TradeSwitch.Buy {
			go BuyBot(analyzeTick)
			continue
		}
		if IsSellFirstPoint(analyzeTick, cond) && global.TradeSwitch.SellFirst {
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
