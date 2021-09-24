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
	fetchentiretick.FetchEntireTick(targetArr, global.LastTradeDayArr, global.TickAnalyzeCondition)
	logger.Logger.Info("Fetch done")
	storeAllEntireTick(targetArr)
	logger.Logger.Info("Begin Training")
	historyCount := getBestHistoryCount(targetArr)
	outSum := getBestOutSum(targetArr, historyCount)
	outInRatio := getBestOutInRatio(targetArr, historyCount, outSum)
	rsiArr := getBestRSI(targetArr, historyCount, outSum, outInRatio)

	global.HistoryCloseCount = int(historyCount)
	cond := global.AnalyzeCondition{
		OutSum:               outSum,
		OutInRatio:           outInRatio,
		CloseDiff:            0,
		CloseChangeRatioLow:  0,
		CloseChangeRatioHigh: 5,
		OpenChangeRatio:      5,
		RsiHigh:              rsiArr[1],
		RsiLow:               rsiArr[0],
	}
	logger.Logger.Info("Finish Training")
	logger.Logger.Warnf("HistoryCount: %d, Cond: %v", global.HistoryCloseCount, cond)
	GetBalance(SearchBuyPoint(targetArr, cond), cond, false)
	time.Sleep(10 * time.Second)
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
func GetBalance(analyzeMap map[string]analyzeentiretick.AnalyzeEntireTick, cond global.AnalyzeCondition, training bool) int64 {
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
	if balance >= maxBalance && !training {
		maxBalance = balance
		logger.Logger.Warnf("Total Balance: %d, TradeCount: %d,HistoryCount: %d, Cond: %v", balance, len(analyzeMap), global.HistoryCloseCount, cond)
	}
	return balance
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

func getBestHistoryCount(targetArr []string) int64 {
	type balanceWithIndex struct {
		index   int
		balance int64
	}
	var tmp balanceWithIndex
	for i := 300; i <= 1300; i += 100 {
		global.HistoryCloseCount = i
		cond := global.AnalyzeCondition{
			OutSum:               200,
			OutInRatio:           60,
			CloseDiff:            0,
			CloseChangeRatioLow:  0,
			CloseChangeRatioHigh: 5,
			OpenChangeRatio:      5,
			RsiHigh:              70,
			RsiLow:               30,
		}
		tmpBalance := GetBalance(SearchBuyPoint(targetArr, cond), cond, true)
		if tmp.balance == 0 {
			tmp.balance = tmpBalance
			tmp.index = i
		} else if tmpBalance > tmp.balance {
			tmp.balance = tmpBalance
			tmp.index = i
		}
	}
	logger.Logger.Warnf("Best HistoryCount: %d, balance: %d", tmp.index, tmp.balance)
	return int64(tmp.index)
}

func getBestOutSum(targetArr []string, historyCloseCount int64) int64 {
	type balanceWithIndex struct {
		index   int
		balance int64
	}
	var tmp balanceWithIndex
	for i := 1000; i >= 100; i -= 10 {
		global.HistoryCloseCount = int(historyCloseCount)
		cond := global.AnalyzeCondition{
			OutSum:               int64(i),
			OutInRatio:           60,
			CloseDiff:            0,
			CloseChangeRatioLow:  0,
			CloseChangeRatioHigh: 5,
			OpenChangeRatio:      5,
			RsiHigh:              70,
			RsiLow:               30,
		}
		tmpBalance := GetBalance(SearchBuyPoint(targetArr, cond), cond, true)
		if tmp.balance == 0 {
			tmp.balance = tmpBalance
			tmp.index = i
		} else if tmpBalance > tmp.balance {
			tmp.balance = tmpBalance
			tmp.index = i
		}
	}
	logger.Logger.Warnf("Best OutSum: %d, balance: %d", tmp.index, tmp.balance)
	return int64(tmp.index)
}

func getBestOutInRatio(targetArr []string, historyCloseCount, outSum int64) float64 {
	type balanceWithIndex struct {
		index   int
		balance int64
	}
	var tmp balanceWithIndex
	for i := 95; i >= 5; i -= 5 {
		global.HistoryCloseCount = int(historyCloseCount)
		cond := global.AnalyzeCondition{
			OutSum:               outSum,
			OutInRatio:           float64(i),
			CloseDiff:            0,
			CloseChangeRatioLow:  0,
			CloseChangeRatioHigh: 5,
			OpenChangeRatio:      5,
			RsiHigh:              70,
			RsiLow:               30,
		}
		tmpBalance := GetBalance(SearchBuyPoint(targetArr, cond), cond, true)
		if tmp.balance == 0 {
			tmp.balance = tmpBalance
			tmp.index = i
		} else if tmpBalance > tmp.balance {
			tmp.balance = tmpBalance
			tmp.index = i
		}
	}
	logger.Logger.Warnf("Best OutInRatio: %d, balance: %d", tmp.index, tmp.balance)
	return float64(tmp.index)
}

func getBestRSI(targetArr []string, historyCloseCount, outSum int64, outInRatio float64) []int64 {
	type balanceWithIndex struct {
		index   int
		balance int64
	}
	var tmp balanceWithIndex
	var rsiArr []int64
	for i := 5; i <= 95; i += 5 {
		global.HistoryCloseCount = int(historyCloseCount)
		cond := global.AnalyzeCondition{
			OutSum:               outSum,
			OutInRatio:           outInRatio,
			CloseDiff:            0,
			CloseChangeRatioLow:  0,
			CloseChangeRatioHigh: 5,
			OpenChangeRatio:      5,
			RsiHigh:              100 - int64(i),
			RsiLow:               int64(i),
		}
		tmpBalance := GetBalance(SearchBuyPoint(targetArr, cond), cond, true)
		if tmp.balance == 0 {
			tmp.balance = tmpBalance
			tmp.index = i
		} else if tmpBalance > tmp.balance {
			tmp.balance = tmpBalance
			tmp.index = i
		}
	}
	rsiArr = append(rsiArr, int64(tmp.index), 100-int64(tmp.index))
	logger.Logger.Warnf("Best RSI Low: %d, RSI High: %d, balance: %d", tmp.index, 100-tmp.index, tmp.balance)
	return rsiArr
}
