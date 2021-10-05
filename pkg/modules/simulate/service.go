// Package simulate package simulate
package simulate

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"

	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/analyzeentiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/entiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/choosetarget"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/entiretickprocess"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/fetchentiretick"
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
	fmt.Print("Use global trade date?(y/n): ")
	reader := bufio.NewReader(os.Stdin)
	ans, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	if ans == "n\n" {
		lastLast := time.Date(2021, 9, 23, 0, 0, 0, 0, time.UTC)
		last := time.Date(2021, 9, 24, 0, 0, 0, 0, time.UTC)
		global.LastTradeDayArr = []time.Time{lastLast, last}
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

// SearchBuyPoint SearchBuyPoint
func SearchBuyPoint(targetArr []string, cond global.AnalyzeCondition) []map[string]analyzeentiretick.AnalyzeEntireTick {
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
	sellFirstPointMap := make(map[string]analyzeentiretick.AnalyzeEntireTick)
	allPoint := simulateAnalyzeEntireMap.GetAllTicks()
	for _, v := range allPoint {
		tmp := v.ToAnalyzeStreamTick()
		tickTimeUnix := time.Unix(0, tmp.TimeStamp)
		lastTime := time.Date(tickTimeUnix.Year(), tickTimeUnix.Month(), tickTimeUnix.Day(), 13, 0, 0, 0, time.Local)
		if tradebot.IsBuyPoint(tmp, cond) {
			if _, ok := buyPointMap[v.StockNum]; !ok && tickTimeUnix.Before(lastTime) {
				buyPointMap[v.StockNum] = *v
				continue
			}
		}
		if tradebot.IsSellFirstPoint(tmp, cond) {
			if _, ok := buyPointMap[v.StockNum]; !ok {
				if _, ok := sellFirstPointMap[v.StockNum]; !ok && tickTimeUnix.Before(lastTime) {
					sellFirstPointMap[v.StockNum] = *v
				}
			}
		}
	}
	return []map[string]analyzeentiretick.AnalyzeEntireTick{buyPointMap, sellFirstPointMap}
}

// GetBalance GetBalance
func GetBalance(analyzeMap []map[string]analyzeentiretick.AnalyzeEntireTick, cond global.AnalyzeCondition, training bool, wg *sync.WaitGroup) {
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
	tmp := bestCond{
		cond:    cond,
		balance: balance,
	}
	if training {
		resultChan <- tmp
	} else {
		logger.Logger.Warnf("Total Balance: %d, TradeCount: %d, Cond: %+v", balance, len(analyzeMap[0])+len(analyzeMap[1]), cond)
	}
}

var resultChan chan bestCond

type bestCond struct {
	cond    global.AnalyzeCondition
	balance int64
}

func catchResult(targetArr []string) {
	var tmp bestCond
	var count, timestamp int
	if timestamp == 0 {
		timestamp = int(time.Now().Unix())
	}
	for {
		result, ok := <-resultChan
		if !ok {
			var wg sync.WaitGroup
			finalCond := tmp.cond
			logger.Logger.Info("Below is best result")
			wg.Add(1)
			go GetBalance(SearchBuyPoint(targetArr, finalCond), finalCond, false, &wg)
			wg.Wait()
			finishSimulate <- 0
			break
		}
		count++
		if tmp.balance == 0 {
			tmp = result
		} else if result.balance > tmp.balance {
			tmp = result
		}
		if count%100 == 0 {
			logger.Logger.Warnf("Finished: %d, Time: %d, Best: %d", count, int(time.Now().Unix())-timestamp, tmp.balance)
			timestamp = int(time.Now().Unix())
		}
	}
}

func getBestCond(targetArr []string) {
	var wg sync.WaitGroup
	var conds []global.AnalyzeCondition
	fmt.Print("Use global cond?(y/n): ")
	reader := bufio.NewReader(os.Stdin)
	ans, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	if ans == "y\n" {
		conds = append(conds, global.TickAnalyzeCondition)
	} else {
		// for j := 400; j >= 100; j -= 100 {
		// 	for m := 70; m >= 50; m -= 5 {
		// for u := 5; u <= 45; u += 5 {
		// for i := 20; i <= 70; i += 5 {
		// 	for z := 5; z <= 25; z += 5 {
		for o := 10; o <= 20; o += 5 {
			for p := 1; p <= 4; p++ {
				for v := 40; v <= 50; v += 10 {
					j := 300
					m := 60
					i := 50
					z := 5
					u := 5
					cond := global.AnalyzeCondition{
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
					conds = append(conds, cond)
				}
			}
		}
	}
	// 	}
	// }
	// }
	// 	}
	// }
	logger.Logger.Warnf("Total simulate counts: %d", len(conds))
	resultChan = make(chan bestCond, len(conds))
	go catchResult(targetArr)
	for _, v := range conds {
		wg.Add(1)
		go GetBalance(SearchBuyPoint(targetArr, v), v, true, &wg)
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
