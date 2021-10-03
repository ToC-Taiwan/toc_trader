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
func SearchBuyPoint(targetArr []string, cond global.AnalyzeCondition) map[string]analyzeentiretick.AnalyzeEntireTick {
	var simulateAnalyzeEntireMap entiretickprocess.AnalyzeEntireTickMap
	var wg sync.WaitGroup
	for _, stockNum := range targetArr {
		ticks := allTickMap.getAllTicksByStockNum(stockNum)
		if len(ticks) == 0 {
			logger.Logger.Errorf("%s has no ticks", stockNum)
			continue
		}
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
		if len(ticks) == 0 {
			logger.Logger.Errorf("%s has no ticks", stockNum)
			continue
		}
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
			if !training {
				logger.Logger.Warnf("Balance: %d, Stock: %s, Name: %s, Total Time: %d", sellCost-buyCost, v.StockNum, global.AllStockNameMap.GetName(v.StockNum), (sellTimeStamp[v.StockNum]-v.TimeStamp)/1000/1000/1000)
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
		logger.Logger.Warnf("Total Balance: %d, TradeCount: %d, Cond: %v", balance, len(analyzeMap), cond)
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
		if count%100 == 0 {
			logger.Logger.Warnf("Finished: %d, Time: %d", count, int(time.Now().Unix())-timestamp)
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
	go catchResult(targetArr)
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
		for j := 500; j >= 300; j -= 100 {
			for l := 200; l >= 150; l -= 50 {
				for m := 70; m >= 55; m -= 5 {
					for n := -3; n <= -1; n++ {
						for k := 7; k >= 5; k-- {
							for i := 40; i <= 55; i += 5 {
								for z := 5; z <= 25; z += 5 {
									for o := 10; o <= 20; o += 5 {
										for p := 1; p <= 4; p++ {
											cond := global.AnalyzeCondition{
												HistoryCloseCount:    int64(j),
												OutSum:               int64(l),
												OutInRatio:           float64(m),
												CloseDiff:            0,
												CloseChangeRatioLow:  float64(n),
												CloseChangeRatioHigh: float64(k),
												OpenChangeRatio:      float64(k),
												RsiHigh:              int64(i + z),
												RsiLow:               int64(i),
												TicksPeriodThreshold: float64(o),
												TicksPeriodLimit:     float64(o) * 1.3,
												TicksPeriodCount:     p,
											}
											conds = append(conds, cond)
										}
									}
								}
							}
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
	wg.Wait()
	close(resultChan)
}
