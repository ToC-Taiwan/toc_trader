// Package tradebot package tradebot
package tradebot

import (
	"time"

	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/analyzestreamtick"
)

var (
	tradeInWaitTime  time.Duration = 30 * time.Second
	tradeOutWaitTime time.Duration = 45 * time.Second
)

// BuyAgent BuyAgent
func BuyAgent(ch chan *analyzestreamtick.AnalyzeStreamTick) {
	for {
		analyzeTick := <-ch
		if checkInBuyMap(analyzeTick.StockNum) {
			logger.GetLogger().Infof("%s already buy", analyzeTick.StockNum)
			continue
		}
		if IsBuyPoint(analyzeTick, global.ForwardCond) {
			logger.GetLogger().Infof("%s %s is on buy point, Close: %.2f", analyzeTick.StockNum, global.AllStockNameMap.GetName(analyzeTick.StockNum), analyzeTick.Close)
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
		if checkInSellFirstMap(analyzeTick.StockNum) {
			logger.GetLogger().Infof("%s already sell first", analyzeTick.StockNum)
			continue
		}
		if IsSellFirstPoint(analyzeTick, global.ReverseCond) {
			logger.GetLogger().Infof("%s %s is on sell first point, Close: %.2f", analyzeTick.StockNum, global.AllStockNameMap.GetName(analyzeTick.StockNum), analyzeTick.Close)
			if global.TradeSwitch.SellFirst {
				SellFirstBot(analyzeTick)
			}
		}
	}
}

func checkInBuyMap(stockNum string) bool {
	return BuyOrderMap.CheckStockExist(stockNum) || FilledBuyOrderMap.CheckStockExist(stockNum)
}

func checkInSellFirstMap(stockNum string) bool {
	return SellFirstOrderMap.CheckStockExist(stockNum) || FilledSellFirstOrderMap.CheckStockExist(stockNum)
}
