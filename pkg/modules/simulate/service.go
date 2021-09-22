// Package simulate package simulate
package simulate

import (
	"sync"
	"time"

	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/analyzeentiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/entiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/choosetarget"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/entiretickprocess"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/fetchentiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/importbasic"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tradebot"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
)

// Simulate Simulate
func Simulate() {
	targetArr, err := choosetarget.GetTopTarget(-1)
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	logger.Logger.Infof("Simulate %d stock", len(targetArr))
	if err := choosetarget.UpdateStockCloseMapByDate(targetArr, global.LastTradeDayArr); err != nil {
		logger.Logger.Error(err)
		return
	}
	global.HistoryCloseCount = 800
	cond := global.AnalyzeCondition{
		OutSum:               260,
		OutInRatio:           60,
		CloseDiff:            0,
		CloseChangeRatioLow:  0,
		CloseChangeRatioHigh: 5,
		OpenChangeRatio:      5,
		RsiHigh:              65,
		RsiLow:               35,
	}
	fetchentiretick.FetchEntireTick(targetArr, global.LastTradeDayArr, cond)
	logger.Logger.Info("Fetch done")

	storeAllEntireTick(targetArr)
	// for i := 5; i < 96; i += 5 {
	// 	global.HistoryCloseCount = 800
	// 	cond := global.AnalyzeCondition{
	// 		OutSum:               260,
	// 		OutInRatio:           60,
	// 		CloseDiff:            0,
	// 		CloseChangeRatioLow:  0,
	// 		CloseChangeRatioHigh: 5,
	// 		OpenChangeRatio:      5,
	// 		RsiHigh:              65,
	// 		RsiLow:               35,
	// 	}
	GetBalance(SearchBuyPoint(targetArr, cond), cond)
	// }
}

// SearchBuyPoint SearchBuyPoint
func SearchBuyPoint(targetArr []string, cond global.AnalyzeCondition) map[string]analyzeentiretick.AnalyzeEntireTick {
	var simulateAnalyzeEntireMap entiretickprocess.AnalyzeEntireTickMap
	var wg sync.WaitGroup
	for _, stockNum := range targetArr {
		wg.Add(1)
		ticks := allTickMap.getAllTicksByStockNum(stockNum)
		ch := make(chan *entiretick.EntireTick)
		saveCh := make(chan []*entiretick.EntireTick)

		lastTradeDay, err := importbasic.GetLastTradeDayByDate(global.LastTradeDay.Format(global.ShortTimeLayout))
		if err != nil {
			panic(err)
		}

		lastClose := global.StockCloseByDateMap.GetClose(stockNum, lastTradeDay.Format(global.ShortTimeLayout))
		if lastClose != 0 {
			go entiretickprocess.TickProcess(stockNum, lastClose, cond, ch, &wg, saveCh, true, &simulateAnalyzeEntireMap)
		} else {
			logger.Logger.Warnf("%s has no %s's close", stockNum, global.LastLastTradeDay.Format(global.ShortTimeLayout))
		}
		for _, v := range ticks {
			tmp := v
			ch <- &tmp
		}
		close(saveCh)
		close(ch)
	}
	wg.Wait()
	buyPointMap := make(map[string]analyzeentiretick.AnalyzeEntireTick)
	allPoint := simulateAnalyzeEntireMap.GetAllTicks()
	for _, v := range allPoint {
		tmp := v.ToAnalyzeStreamTick()
		if tradebot.IsBuyPoint(tmp, cond) {
			tickTimeUnix := time.Unix(0, tmp.TimeStamp)
			lastTime := time.Date(tickTimeUnix.Year(), tickTimeUnix.Month(), tickTimeUnix.Day(), 13, 0, 0, 0, time.Local)
			if _, ok := buyPointMap[v.StockNum]; !ok && tickTimeUnix.Before(lastTime) {
				buyPointMap[v.StockNum] = *v
			}
		}
	}
	return buyPointMap
}

var maxBalance int64

// GetBalance GetBalance
func GetBalance(analyzeMap map[string]analyzeentiretick.AnalyzeEntireTick, cond global.AnalyzeCondition) {
	sellTimeStamp := make(map[string]int64)
	var balance int64
	for stockNum, v := range analyzeMap {
		ticks := allTickMap.getAllTicksByStockNum(stockNum)
		var historyClose []float64
		var buyPrice, sellPrice float64
		for _, k := range ticks {
			historyClose = append(historyClose, k.Close)
			if len(historyClose) > global.HistoryCloseCount {
				historyClose = historyClose[1:]
			}
			if k.TimeStamp == v.TimeStamp && buyPrice == 0 {
				historyClose = []float64{}
				buyPrice = k.Close
			}
			if buyPrice != 0 {
				sellPrice = tradebot.GetSellPrice(k.ToStreamTick(), time.Unix(0, v.TimeStamp).Add(-8*time.Hour), historyClose, buyPrice, cond)
				if sellPrice != 0 {
					sellTimeStamp[k.StockNum] = k.TimeStamp
					break
				}
			}
		}
		if sellPrice == 0 {
			logger.Logger.Warnf("%s no sell point", stockNum)
		} else {
			buyCost := tradebot.GetStockBuyCost(buyPrice, global.OneTimeQuantity)
			sellCost := tradebot.GetStockSellCost(sellPrice, global.OneTimeQuantity)
			balance += (sellCost - buyCost)
			logger.Logger.Warnf("Balance: %d, Stock: %s, Name: %s, Total Time: %d", sellCost-buyCost, v.StockNum, global.AllStockNameMap.GetName(v.StockNum), (sellTimeStamp[v.StockNum]-v.TimeStamp)/1000/1000/1000)
		}
	}
	if balance >= maxBalance {
		maxBalance = balance
		logger.Logger.Warnf("Total Balance: %d, TradeCount: %d,HistoryCount: %d, Cond: %v", balance, len(analyzeMap), global.HistoryCloseCount, cond)
	}
}

var allTickMap entireTickMap

func storeAllEntireTick(stockArr []string) {
	for _, stockNum := range stockArr {
		ticks, err := entiretick.GetAllEntiretickByStock(stockNum, global.GlobalDB)
		if err != nil {
			logger.Logger.Error(err)
			continue
		}
		allTickMap.saveByStockNum(stockNum, ticks)
	}
}
