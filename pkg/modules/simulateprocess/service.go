// Package simulateprocess package simulateprocess
package simulateprocess

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/analyzeentiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/entiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulate"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/choosetarget"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/entiretickprocess"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/fetchentiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/importbasic"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tradebot"
	"gitlab.tocraw.com/root/toc_trader/tools/common"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
)

var balanceType string
var allTickMap entireTickMap
var resultChan chan simulate.Result
var totalTimesChan chan int
var totalTimes, finishTimes int
var simulateDayArr []time.Time
var targetArrMap targetArrMutex
var discardOverTime bool

// Simulate Simulate
func Simulate() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("* Balance type?(a: forward, b: reverse, c: force_both): ")
	balanceTypeAns, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	balanceTypeAns = strings.TrimSuffix(balanceTypeAns, "\n")
	switch balanceTypeAns {
	case "a":
		balanceType = "a"
	case "b":
		balanceType = "b"
	case "c":
		balanceType = "c"
	}

	fmt.Print("* Discard over time trade?(y/n): ")
	discardOverTimeAns, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	if discardOverTimeAns = strings.TrimSuffix(discardOverTimeAns, "\n"); discardOverTimeAns == "y\n" {
		discardOverTime = true
	}

	var useGlobal bool
	fmt.Print("* Use global cond?(y/n): ")
	useGlobalAns, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	if useGlobalAns = strings.TrimSuffix(useGlobalAns, "\n"); useGlobalAns == "y\n" {
		useGlobal = true
	}

	fmt.Print("* N days?: ")
	nDaysAns, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	nDaysAns = strings.TrimSuffix(nDaysAns, "\n")
	n, err := common.StrToInt64(nDaysAns)
	if err != nil {
		panic(err)
	}
	clearAllSimulation()
	tradeDayArr, err := importbasic.GetLastNTradeDay(n + 1)
	if err != nil {
		panic(err)
	}

	for i, date := range tradeDayArr {
		if i == 0 {
			continue
		}
		if targets, err := choosetarget.GetTargetByVolumeRankByDate(date.Format(global.ShortTimeLayout), 200); err != nil {
			panic(err)
		} else {
			for i, v := range targets {
				fmt.Printf("%s volume rank no. %d is %s\n", date.Format(global.ShortTimeLayout), i+1, global.AllStockNameMap.GetName(v))
			}
			targetArrMap.saveByDate(tradeDayArr[i-1].Format(global.ShortTimeLayout), targets)
			for {
				tmp := []time.Time{date}
				err := choosetarget.UpdateStockCloseMapByDate(targets, tmp)
				if err != nil {
					logger.Logger.Error(err)
				} else {
					break
				}
			}
			tmp := []time.Time{tradeDayArr[i-1]}
			fetchentiretick.FetchEntireTick(targets, tmp, global.TickAnalyzeCondition)
			storeAllEntireTick(targets, tmp)
		}
	}
	simulateDayArr = tradeDayArr
	logger.Logger.Info("Fetch done")

	resultChan = make(chan simulate.Result)
	totalTimesChan = make(chan int)
	go totalTimesReceiver()
	go catchResult(useGlobal)

	var wg sync.WaitGroup
	if useGlobal {
		getBestCond(int(global.TickAnalyzeCondition.HistoryCloseCount), useGlobal)
	} else {
		for i := 2500; i >= 1500; i -= 500 {
			wg.Add(1)
			go func(historyCount int) {
				defer wg.Done()
				getBestCond(historyCount, useGlobal)
			}(i)
		}
	}
	wg.Wait()

	close(resultChan)
	logger.Logger.Warn("Finish simulate")
	time.Sleep(10 * time.Second)
}

// getBestCond getBestCond
func getBestCond(historyCount int, useGlobal bool) {
	var wg sync.WaitGroup
	var conds []*simulationcond.AnalyzeCondition
	if useGlobal {
		conds = append(conds, &global.TickAnalyzeCondition)
	} else {
		for m := 55; m >= 55; m -= 5 {
			for u := 5; u <= 15; u += 5 {
				for i := 50; i >= 50; i -= 5 {
					for z := 0; z <= 5; z++ {
						for o := 11; o >= 1; o -= 2 {
							for p := 1; p <= 5; p++ {
								for v := 3; v <= 6; v++ {
									cond := simulationcond.AnalyzeCondition{
										HistoryCloseCount:    int64(historyCount),
										OutInRatio:           float64(m),
										ReverseOutInRatio:    float64(u),
										CloseDiff:            0,
										CloseChangeRatioLow:  -1,
										CloseChangeRatioHigh: 8,
										OpenChangeRatio:      4,
										RsiHigh:              float64(i) + float64(z)/10,
										RsiLow:               float64(i),
										ReverseRsiHigh:       float64(i) + float64(z)/10,
										ReverseRsiLow:        float64(i),
										TicksPeriodThreshold: float64(o),
										TicksPeriodLimit:     float64(o) * 1.3,
										TicksPeriodCount:     p,
										VolumePerSecond:      int64(v),
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
		go GetBalance(SearchTradePoint(simulateDayArr, *v), *v, training, &wg)
	}
	wg.Wait()
}

// SearchTradePoint SearchTradePoint
func SearchTradePoint(tradeDayArr []time.Time, cond simulationcond.AnalyzeCondition) (pointMapArr map[string][]map[string]*analyzeentiretick.AnalyzeEntireTick) {
	pointMapArr = make(map[string][]map[string]*analyzeentiretick.AnalyzeEntireTick)
	var simulateAnalyzeEntireMap entiretickprocess.AnalyzeEntireTickMap
	for i, date := range tradeDayArr {
		simulateAnalyzeEntireMap.ClearAll()
		if i == len(tradeDayArr)-1 {
			break
		}
		var wg sync.WaitGroup
		targetArr := targetArrMap.getArrByDate(date.Format(global.ShortTimeLayout))
		for _, stockNum := range targetArr {
			ticks := allTickMap.getAllTicksByStockNumAndDate(stockNum, tradeDayArr[i].Format(global.ShortTimeLayout))
			wg.Add(1)
			ch := make(chan *entiretick.EntireTick)
			saveCh := make(chan []*entiretick.EntireTick)
			lastClose := global.StockCloseByDateMap.GetClose(stockNum, tradeDayArr[i+1].Format(global.ShortTimeLayout))
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
		pointMapArr[date.Format(global.ShortTimeLayout)] = append(pointMapArr[date.Format(global.ShortTimeLayout)], buyPointMap, sellFirstPointMap)
	}
	return pointMapArr
}

// GetBalance GetBalance
func GetBalance(analyzeMapMap map[string][]map[string]*analyzeentiretick.AnalyzeEntireTick, cond simulationcond.AnalyzeCondition, training bool, wg *sync.WaitGroup) {
	defer wg.Done()
	var forwardBalance, reverseBalance int64
	var tradeCount int64
	for date, analyzeMap := range analyzeMapMap {
		if balanceType == "c" && (len(analyzeMap[0]) == 0 || len(analyzeMap[1]) == 0) && training {
			totalTimesChan <- -1
			return
		}
		sellTimeStamp := make(map[string]int64)
		for stockNum, v := range analyzeMap[0] {
			ticks := allTickMap.getAllTicksByStockNumAndDate(stockNum, date)
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
				back := tradebot.GetStockTradeFeeDiscount(buyPrice, global.OneTimeQuantity) + tradebot.GetStockTradeFeeDiscount(sellPrice, global.OneTimeQuantity)
				forwardBalance += (sellCost - buyCost + back)
				if sellTimeStamp[v.StockNum] > endTradeTime && training && discardOverTime {
					totalTimesChan <- -1
					return
				}
				tradeCount++
				if !training {
					logger.Logger.Warnf("%s Forward Balance: %d, Stock: %s, Name: %s, Total Time: %d, %.2f, %.2f", date, sellCost-buyCost+back, v.StockNum, global.AllStockNameMap.GetName(v.StockNum), (sellTimeStamp[v.StockNum]-v.TimeStamp)/1000/1000/1000, buyPrice, sellPrice)
				}
			}
		}
		sellTimeStamp = make(map[string]int64)
		for stockNum, v := range analyzeMap[1] {
			ticks := allTickMap.getAllTicksByStockNumAndDate(stockNum, date)
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
				back := tradebot.GetStockTradeFeeDiscount(buyLaterPrice, global.OneTimeQuantity) + tradebot.GetStockTradeFeeDiscount(sellFirstPrice, global.OneTimeQuantity)
				if sellTimeStamp[v.StockNum] > endTradeTime && training && discardOverTime {
					totalTimesChan <- -1
					return
				}
				tradeCount++
				reverseBalance += (sellCost - buyCost + back)
				if !training {
					logger.Logger.Warnf("%s Reverse Balance: %d, Stock: %s, Name: %s, Total Time: %d, %.2f, %.2f", date, sellCost-buyCost+back, v.StockNum, global.AllStockNameMap.GetName(v.StockNum), (sellTimeStamp[v.StockNum]-v.TimeStamp)/1000/1000/1000, buyLaterPrice, sellFirstPrice)
				}
			}
		}
	}
	tmp := simulate.Result{
		Balance:        forwardBalance + reverseBalance,
		ForwardBalance: forwardBalance,
		ReverseBalance: reverseBalance,
		TradeCount:     tradeCount,
		Cond:           cond,
	}
	if training {
		resultChan <- tmp
	} else {
		logger.Logger.Warnf("Total Balance: %d, TradeCount: %d", tmp.Balance, tradeCount)
		logger.Logger.Warnf("Cond: %+v", cond)
	}
}

func catchResult(useGlobal bool) {
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
					go GetBalance(SearchTradePoint(simulateDayArr, k.Cond), k.Cond, false, &wg)
				}
				wg.Wait()
			} else if !useGlobal {
				logger.Logger.Info("No best result")
			}
			break
		}
		count++
		totalTimesChan <- -1
		switch {
		case len(tmp) == 0:
			tmp = append(tmp, result)
		case result.Balance > tmp[0].Balance:
			tmp = []simulate.Result{}
			tmp = append(tmp, result)
		case result.Balance == tmp[0].Balance:
			tmp = append(tmp, result)
		}
		if count%10 == 0 {
			// logger.Logger.Warnf("Best: %d", tmp[0].Balance)
			if err := simulate.InsertMultiRecord(save, global.GlobalDB); err != nil {
				logger.Logger.Error(err)
			}
			save = []simulate.Result{}
		}
	}
}

func totalTimesReceiver() {
	var count int
	for {
		times := <-totalTimesChan
		if times > 0 {
			count++
			totalTimes += times
			if count == 3 {
				go progressBar(totalTimes)
			}
			// logger.Logger.Warnf("Total simulate counts: %d", totalTimes)
		} else {
			finishTimes++
			if count == 3 {
				bar.Increment()
			}
		}
	}
}

var bar *mpb.Bar

func progressBar(total int) {
	p := mpb.New(mpb.WithWidth(64))
	name := "Time Left:"
	bar = p.AddBar(int64(total),
		mpb.PrependDecorators(
			decor.Name(name, decor.WC{W: len(name) + 1, C: decor.DidentRight}),
			decor.OnComplete(
				decor.AverageETA(decor.ET_STYLE_GO, decor.WC{W: 4}), "done",
			),
		),
		// mpb.AppendDecorators(decor.Percentage()),
		mpb.AppendDecorators(decor.Counters(0, "")),
	)
	p.Wait()
}

func storeAllEntireTick(stockArr []string, tradeDayArr []time.Time) {
	for _, stockNum := range stockArr {
		for _, date := range tradeDayArr {
			ticks, err := entiretick.GetAllEntiretickByStockByDate(stockNum, date.Format(global.ShortTimeLayout), global.GlobalDB)
			if err != nil {
				logger.Logger.Error(err)
				continue
			}
			allTickMap.saveByStockNumAndDate(stockNum, date.Format(global.ShortTimeLayout), ticks)
		}
	}
}

func clearAllSimulation() {
	if err := simulate.DeleteAll(global.GlobalDB); err != nil {
		panic(err)
	}
	if err := simulationcond.DeleteAll(global.GlobalDB); err != nil {
		panic(err)
	}
}
