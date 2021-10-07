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

var finishSimulate chan int

var balanceType string

// Simulate Simulate
func Simulate() {
	if err := simulate.DeleteAll(global.GlobalDB); err != nil {
		panic(err)
	}
	finishSimulate = make(chan int)
	targetArr, err := choosetarget.GetTopTarget(-1)
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	// fmt.Print("Use global trade date?(y/n): ")
	// reader := bufio.NewReader(os.Stdin)
	// ans, err := reader.ReadString('\n')
	// if err != nil {
	// 	panic(err)
	// }
	// if ans == "n\n" {
	// 	lastLast := time.Date(2021, 9, 23, 0, 0, 0, 0, time.UTC)
	// 	last := time.Date(2021, 9, 24, 0, 0, 0, 0, time.UTC)
	// 	global.LastTradeDayArr = []time.Time{lastLast, last}
	// }
	fmt.Print("Simulate balance type?(a: forward, b: reverse, c: all): ")
	reader := bufio.NewReader(os.Stdin)
	ans2, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	switch ans2 {
	case "a\n":
		balanceType = "a"
	case "b\n":
		balanceType = "b"
	case "c\n":
		balanceType = "c"
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
			time.Sleep(10 * time.Second)
			return
		}
	}
}

// SearchTradePoint SearchTradePoint
func SearchTradePoint(targetArr []string, cond simulationcond.AnalyzeCondition) []map[string]*analyzeentiretick.AnalyzeEntireTick {
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
		if tradebot.IsBuyPoint(tmp, cond) && (balanceType == "a" || balanceType == "c") {
			if _, ok := buyPointMap[v.StockNum]; !ok && tickTimeUnix.Before(lastTime) {
				buyPointMap[v.StockNum] = v
				continue
			}
		}
		if tradebot.IsSellFirstPoint(tmp, cond) && (balanceType == "b" || balanceType == "c") {
			if _, ok := buyPointMap[v.StockNum]; !ok {
				if _, ok := sellFirstPointMap[v.StockNum]; !ok && tickTimeUnix.Before(lastTime) {
					sellFirstPointMap[v.StockNum] = v
				}
			}
		}
	}
	return []map[string]*analyzeentiretick.AnalyzeEntireTick{buyPointMap, sellFirstPointMap}
}

// GetBalance GetBalance
func GetBalance(analyzeMap []map[string]*analyzeentiretick.AnalyzeEntireTick, cond simulationcond.AnalyzeCondition, training bool, wg *sync.WaitGroup) {
	defer wg.Done()
	sellTimeStamp := make(map[string]int64)
	var balance int64
	for stockNum, v := range analyzeMap[0] {
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
			if (sellTimeStamp[v.StockNum]-v.TimeStamp)/1000/1000/1000 > 10800 {
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
			if (sellTimeStamp[v.StockNum]-v.TimeStamp)/1000/1000/1000 > 10800 {
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

var resultChan chan simulate.Result

// type bestCond struct {
// 	cond    simulationcond.AnalyzeCondition
// 	balance int64
// }

func catchResult(targetArr []string, times int) {
	var save []simulate.Result
	var tmp []simulate.Result
	var consumeArr []int
	var count int
	timestamp := int(time.Now().Unix())
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
			// finalCond := tmp.cond
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
			finishSimulate <- 0
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
			var total, average float64
			consume := int(time.Now().Unix()) - timestamp
			consumeArr = append(consumeArr, consume)
			for _, v := range consumeArr {
				total += float64(v)
			}
			average = total / float64(len(consumeArr))
			rareTime := average * (float64(times-count) / 100) / 60
			logger.Logger.Warnf("Finished: %d, Time: %d, Estimate: %.2f min, Best: %d", count, consume, rareTime, tmp[0].Balance)
			timestamp = int(time.Now().Unix())
			if err := simulate.InsertMultiRecord(save, global.GlobalDB); err != nil {
				logger.Logger.Error(err)
			}
			save = []simulate.Result{}
		}
	}
}

func getBestCond(targetArr []string) {
	var wg sync.WaitGroup
	var conds []*simulationcond.AnalyzeCondition
	fmt.Print("Use global cond?(y/n): ")
	reader := bufio.NewReader(os.Stdin)
	ans, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	if ans == "y\n" {
		conds = append(conds, &global.TickAnalyzeCondition)
	} else {
		for j := 3000; j >= 3000; j -= 500 {
			for m := 60; m >= 50; m -= 5 {
				for u := 10; u <= 20; u += 5 {
					for i := 45; i <= 50; i += 5 {
						for z := 0; z <= 10; z += 5 {
							for o := 20; o >= 5; o -= 5 {
								for p := 3; p >= 1; p-- {
									for v := 200; v >= 100; v -= 10 {
										// j := 1000
										// m := 55
										// u := 25
										// i := 50
										// z := 0
										// o := 10
										// p := 2
										// v := 180
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
	}
	if err := simulationcond.DeleteAll(global.GlobalDB); err != nil {
		panic(err)
	}
	if err := simulationcond.InsertMultiRecord(conds, global.GlobalDB); err != nil {
		panic(err)
	}
	logger.Logger.Warnf("Total simulate counts: %d", len(conds))
	resultChan = make(chan simulate.Result)
	go catchResult(targetArr, len(conds))
	for _, v := range conds {
		wg.Add(1)
		go GetBalance(SearchTradePoint(targetArr, *v), *v, true, &wg)
	}
	wg.Wait()
	close(resultChan)
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
