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

var finishSimulate chan int

// Simulate Simulate
func Simulate() {
	finishSimulate = make(chan int)
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
	getBestCond(targetArr)
	for {
		finish := <-finishSimulate
		if finish == 0 {
			close(finishSimulate)
			return
		}
	}
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

// GetBalance GetBalance
func GetBalance(analyzeMap map[string]analyzeentiretick.AnalyzeEntireTick, cond global.AnalyzeCondition, training bool, wg *sync.WaitGroup) {
	defer wg.Done()
	sellTimeStamp := make(map[string]int64)
	var balance int64
	for stockNum, v := range analyzeMap {
		ticks := allTickMap.getAllTicksByStockNum(stockNum)
		var historyClose []float64
		var buyPrice, sellPrice float64
		for _, k := range ticks {
			historyClose = append(historyClose, k.Close)
			if len(historyClose) > int(cond.HistoryCloseCount) {
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
		}
	}
	tmp := bestCond{
		historyCount:    int(cond.HistoryCloseCount),
		outSum:          int(cond.OutSum),
		outInRatio:      int(cond.OutInRatio),
		closeLow:        int(cond.CloseChangeRatioLow),
		closeHigh:       int(cond.CloseChangeRatioHigh),
		openChangeRatio: int(cond.OpenChangeRatio),
		rsiLow:          int(cond.RsiLow),
		rsiHigh:         int(cond.RsiHigh),
		balance:         balance,
	}
	resultChan <- tmp
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

var resultChan chan bestCond

type bestCond struct {
	historyCount    int
	outSum          int
	outInRatio      int
	closeLow        int
	closeHigh       int
	openChangeRatio int
	rsiLow          int
	rsiHigh         int
	balance         int64
}

func catchResult() {
	var tmp bestCond
	var count, timestamp int
	for {
		result, ok := <-resultChan
		if !ok {
			logger.Logger.Warnf("Best: %+v", tmp)
			finalCond := global.AnalyzeCondition{
				HistoryCloseCount:    int64(result.historyCount),
				OutSum:               int64(result.outSum),
				OutInRatio:           float64(result.outInRatio),
				CloseDiff:            0,
				CloseChangeRatioLow:  float64(result.closeLow),
				CloseChangeRatioHigh: float64(result.closeHigh),
				OpenChangeRatio:      float64(result.openChangeRatio),
				RsiHigh:              int64(result.rsiHigh),
				RsiLow:               int64(result.rsiLow),
			}
			global.TickAnalyzeCondition = finalCond
			finishSimulate <- 0
			break
		}
		count++
		if count%100 == 0 {
			logger.Logger.Warn(count, int(time.Now().Unix())-timestamp)
			timestamp = int(time.Now().Unix())
		}
		if tmp.balance == 0 {
			tmp = result
		} else if result.balance > tmp.balance {
			tmp = result
		}
	}
}

func getBestCond(targetArr []string) {
	resultChan = make(chan bestCond)
	go catchResult()
	var wg sync.WaitGroup
	var conds []global.AnalyzeCondition
	for j := 200; j <= 500; j += 100 {
		for l := 200; l >= 100; l -= 100 {
			for m := 75; m >= 50; m -= 5 {
				for n := -3; n <= 1; n++ {
					for k := 10; k >= 4; k-- {
						for i := 40; i <= 60; i += 5 {
							cond := global.AnalyzeCondition{
								HistoryCloseCount:    int64(j),
								OutSum:               int64(l),
								OutInRatio:           float64(m),
								CloseChangeRatioLow:  float64(n),
								CloseChangeRatioHigh: float64(k),
								OpenChangeRatio:      float64(k),
								RsiLow:               int64(i),
								RsiHigh:              int64(i),
								CloseDiff:            0,
							}
							conds = append(conds, cond)
							// wg.Add(1)
							// go GetBalance(SearchBuyPoint(targetArr, cond), cond, true, &wg)
						}
					}
				}
			}
		}
	}
	logger.Logger.Warnf("Total simulate counts: %d", len(conds))
	for _, v := range conds {
		wg.Add(1)
		go GetBalance(SearchBuyPoint(targetArr, v), v, true, &wg)
	}
	// cond := global.AnalyzeCondition{
	// 	HistoryCloseCount:    200,
	// 	OutSum:               100,
	// 	OutInRatio:           40,
	// 	CloseDiff:            0,
	// 	CloseChangeRatioLow:  -10,
	// 	CloseChangeRatioHigh: 10,
	// 	OpenChangeRatio:      10,
	// 	RsiHigh:              55,
	// 	RsiLow:               55,
	// }
	// wg.Add(1)
	// go GetBalance(SearchBuyPoint(targetArr, cond), cond, true, &wg)
	wg.Wait()
	close(resultChan)
}
