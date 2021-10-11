// Package simulateprocess package simulateprocess
package simulateprocess

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"

	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/analyzeentiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/entiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulate"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/choosetarget"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/entiretickprocess"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/fetchentiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tradebot"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
)

var balanceType string
var allTickMap entireTickMap
var resultChan chan simulate.Result
var totalTimesChan chan int
var totalTimes int

// Simulate Simulate
func Simulate() {
	if err := simulate.DeleteAll(global.GlobalDB); err != nil {
		panic(err)
	}
	if err := simulationcond.DeleteAll(global.GlobalDB); err != nil {
		panic(err)
	}
	targetArr, err := choosetarget.GetTopTarget(-1)
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	fmt.Print("Simulate balance type?(a: forward, b: reverse, c: force_both): ")
	reader := bufio.NewReader(os.Stdin)
	ans, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	switch ans {
	case "a\n":
		balanceType = "a"
	case "b\n":
		balanceType = "b"
	case "c\n":
		balanceType = "c"
	}
	var useGlobal bool
	fmt.Print("Use global cond?(y/n): ")
	ans2, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	if ans2 == "y\n" {
		useGlobal = true
	}
	logger.Logger.Infof("Simulate %d stock", len(targetArr))
	if err := choosetarget.UpdateStockCloseMapByDate(targetArr, global.LastTradeDayArr); err != nil {
		logger.Logger.Error(err)
		return
	}
	fetchentiretick.FetchEntireTick(targetArr, global.LastTradeDayArr, global.TickAnalyzeCondition)
	logger.Logger.Info("Fetch done")
	storeAllEntireTick(targetArr)
	resultChan = make(chan simulate.Result)
	totalTimesChan = make(chan int)
	go catchResult(targetArr)
	go totalTimesReceiver()
	var wg sync.WaitGroup
	if useGlobal {
		getBestCond(targetArr, int(global.TickAnalyzeCondition.HistoryCloseCount), useGlobal)
	} else {
		for i := 2500; i >= 500; i -= 500 {
			wg.Add(1)
			go func(historyCount int) {
				defer wg.Done()
				getBestCond(targetArr, historyCount, useGlobal)
			}(i)
		}
	}
	wg.Wait()
	close(resultChan)
	logger.Logger.Warn("Finish simulate")
	time.Sleep(10 * time.Second)
}

// SearchTradePoint SearchTradePoint
func SearchTradePoint(targetArr []string, cond simulationcond.AnalyzeCondition) (pointMapArr []map[string]*analyzeentiretick.AnalyzeEntireTick) {
	var simulateAnalyzeEntireMap entiretickprocess.AnalyzeEntireTickMap
	var wg sync.WaitGroup
	for _, stockNum := range targetArr {
		ticks := allTickMap.getAllTicksByStockNum(stockNum)
		wg.Add(1)
		ch := make(chan *entiretick.EntireTick)
		saveCh := make(chan []*entiretick.EntireTick)
		lastClose := global.StockCloseByDateMap.GetClose(stockNum, global.LastTradeDayArr[0].Format(global.ShortTimeLayout))
		if lastClose != 0 {
			go entiretickprocess.TickProcess(stockNum, lastClose, cond, ch, &wg, saveCh, true, &simulateAnalyzeEntireMap)
		} else {
			logger.Logger.Warnf("%s has no %s's close", stockNum, global.LastLastTradeDay.Format(global.ShortTimeLayout))
			continue
		}
		for _, v := range ticks {
			tmp := v
			ch <- tmp
		}
		close(saveCh)
		close(ch)
	}
	wg.Wait()
	buyPointMap := make(map[string]*analyzeentiretick.AnalyzeEntireTick)
	sellFirstPointMap := make(map[string]*analyzeentiretick.AnalyzeEntireTick)
	allPoint := simulateAnalyzeEntireMap.GetAllTicks()
	for _, v := range allPoint {
		tmp := v.ToAnalyzeStreamTick()
		tickTimeUnix := time.Unix(0, tmp.TimeStamp)
		lastTime := time.Date(tickTimeUnix.Year(), tickTimeUnix.Month(), tickTimeUnix.Day(), 13, 0, 0, 0, time.Local)
		if tickTimeUnix.After(lastTime) || buyPointMap[v.StockNum] != nil || sellFirstPointMap[v.StockNum] != nil {
			continue
		}
		if tradebot.IsBuyPoint(tmp, cond) && (balanceType == "a" || balanceType == "c") {
			buyPointMap[v.StockNum] = v
		} else if tradebot.IsSellFirstPoint(tmp, cond) && (balanceType == "b" || balanceType == "c") {
			sellFirstPointMap[v.StockNum] = v
		}
	}
	pointMapArr = append(pointMapArr, buyPointMap, sellFirstPointMap)
	return pointMapArr
}

// GetBalance GetBalance
func GetBalance(analyzeMap []map[string]*analyzeentiretick.AnalyzeEntireTick, cond simulationcond.AnalyzeCondition, training bool, wg *sync.WaitGroup) {
	defer wg.Done()
	if balanceType == "c" && (len(analyzeMap[0]) == 0 || len(analyzeMap[1]) == 0) {
		totalTimesChan <- -1
		return
	}
	sellTimeStamp := make(map[string]int64)
	var balance int64
	for stockNum, v := range analyzeMap[0] {
		ticks := allTickMap.getAllTicksByStockNum(stockNum)
		endTradeTime := getLastTradeTimeByEntireTickTimeStamp(ticks[0].TimeStamp)
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
			if sellTimeStamp[v.StockNum]-endTradeTime > 10800*1000*1000*1000 && training {
				return
			}
			if !training {
				logger.Logger.Warnf("Forward Balance: %d, Stock: %s, Name: %s, Total Time: %d, %.2f, %.2f", sellCost-buyCost, v.StockNum, global.AllStockNameMap.GetName(v.StockNum), (sellTimeStamp[v.StockNum]-v.TimeStamp)/1000/1000/1000, buyPrice, sellPrice)
			}
		}
	}
	sellTimeStamp = make(map[string]int64)
	for stockNum, v := range analyzeMap[1] {
		ticks := allTickMap.getAllTicksByStockNum(stockNum)
		endTradeTime := getLastTradeTimeByEntireTickTimeStamp(ticks[0].TimeStamp)
		var historyClose []float64
		var sellFirstPrice, buyLaterPrice float64
		for _, k := range ticks {
			historyClose = append(historyClose, k.Close)
			if len(historyClose) > int(cond.HistoryCloseCount) {
				historyClose = historyClose[1:]
			}
			if k.TimeStamp == v.TimeStamp && sellFirstPrice == 0 {
				historyClose = []float64{}
				sellFirstPrice = k.Close
			}
			if sellFirstPrice != 0 {
				buyLaterPrice = tradebot.GetBuyLaterPrice(k.ToStreamTick(), time.Unix(0, v.TimeStamp).Add(-8*time.Hour), historyClose, sellFirstPrice, cond)
				if buyLaterPrice != 0 {
					sellTimeStamp[k.StockNum] = k.TimeStamp
					break
				}
			}
		}
		if buyLaterPrice == 0 {
			logger.Logger.Warnf("%s no buy later point", stockNum)
		} else {
			buyCost := tradebot.GetStockBuyCost(buyLaterPrice, global.OneTimeQuantity)
			sellCost := tradebot.GetStockSellCost(sellFirstPrice, global.OneTimeQuantity)
			if sellTimeStamp[v.StockNum]-endTradeTime > 10800*1000*1000*1000 && training {
				return
			}
			balance += (sellCost - buyCost)
			if !training {
				logger.Logger.Warnf("Reverse Balance: %d, Stock: %s, Name: %s, Total Time: %d, %.2f, %.2f", sellCost-buyCost, v.StockNum, global.AllStockNameMap.GetName(v.StockNum), (sellTimeStamp[v.StockNum]-v.TimeStamp)/1000/1000/1000, buyLaterPrice, sellFirstPrice)
			}
		}
	}
	tmp := simulate.Result{
		Cond:    cond,
		Balance: balance,
	}
	if training {
		resultChan <- tmp
	} else {
		logger.Logger.Warnf("Total Balance: %d, TradeCount: %d, Cond: %+v", balance, len(analyzeMap[0])+len(analyzeMap[1]), cond)
	}
}

func catchResult(targetArr []string) {
	var save []simulate.Result
	var tmp []simulate.Result
	var count int
	for {
		result, ok := <-resultChan
		if result.Cond.Model.ID != 0 {
			save = append(save, result)
		}
		if !ok {
			if err := simulate.InsertMultiRecord(save, global.GlobalDB); err != nil {
				logger.Logger.Error(err)
			}
			var wg sync.WaitGroup
			if len(tmp) > 0 {
				logger.Logger.Info("Below is best result")
				for _, k := range tmp {
					wg.Add(1)
					go GetBalance(SearchTradePoint(targetArr, k.Cond), k.Cond, false, &wg)
				}
				wg.Wait()
			} else {
				logger.Logger.Info("No best result")
			}
			break
		}
		count++
		switch {
		case len(tmp) == 0:
			tmp = append(tmp, result)
		case result.Balance > tmp[0].Balance:
			tmp = []simulate.Result{}
			tmp = append(tmp, result)
		case result.Balance == tmp[0].Balance:
			tmp = append(tmp, result)
		}
		if count%100 == 0 {
			logger.Logger.Warnf("Finished: %d, Rare: %d, Best: %d", count, totalTimes-count, tmp[0].Balance)
			if err := simulate.InsertMultiRecord(save, global.GlobalDB); err != nil {
				logger.Logger.Error(err)
			}
			save = []simulate.Result{}
		}
	}
}

func getBestCond(targetArr []string, historyCount int, useGlobal bool) {
	var wg sync.WaitGroup
	var conds []*simulationcond.AnalyzeCondition
	if useGlobal {
		conds = append(conds, &global.TickAnalyzeCondition)
	} else {
		for m := 75; m >= 50; m -= 5 {
			for u := 5; u <= 25; u += 5 {
				for i := 35; i <= 50; i += 5 {
					for z := 0; z <= 15; z += 5 {
						for o := 20; o >= 5; o -= 5 {
							for p := 3; p >= 1; p-- {
								for v := 250; v >= 50; v -= 50 {
									j := historyCount
									cond := simulationcond.AnalyzeCondition{
										HistoryCloseCount:    int64(j),
										OutInRatio:           float64(m),
										ReverseOutInRatio:    float64(u),
										CloseDiff:            0,
										CloseChangeRatioLow:  -3,
										CloseChangeRatioHigh: 6,
										OpenChangeRatio:      6,
										RsiHigh:              int64(i + z),
										RsiLow:               int64(i),
										ReverseRsiHigh:       int64(i + z),
										ReverseRsiLow:        int64(i),
										TicksPeriodThreshold: float64(o),
										TicksPeriodLimit:     float64(o) * 1.3,
										TicksPeriodCount:     p,
										Volume:               int64(v),
									}
									conds = append(conds, &cond)
								}
							}
						}
					}
				}
			}
		}
	}
	if err := simulationcond.InsertMultiRecord(conds, global.GlobalDB); err != nil {
		panic(err)
	}
	totalTimesChan <- len(conds)
	training := true
	if useGlobal {
		training = false
	}
	for _, v := range conds {
		wg.Add(1)
		go GetBalance(SearchTradePoint(targetArr, *v), *v, training, &wg)
	}
	wg.Wait()
}

func totalTimesReceiver() {
	for {
		times := <-totalTimesChan
		totalTimes += times
		if times > 0 {
			logger.Logger.Warnf("Total simulate counts: %d", totalTimes)
		}
	}
}

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
