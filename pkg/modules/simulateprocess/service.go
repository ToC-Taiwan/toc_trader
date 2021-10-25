// Package simulateprocess package simulateprocess
package simulateprocess

import (
	"fmt"
	"sync"
	"time"

	"github.com/manifoldco/promptui"
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
	"gitlab.tocraw.com/root/toc_trader/tools/db"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
)

var (
	balanceType             string
	allTickMap              entireTickMap
	resultChan              chan simulate.Result
	totalTimesChan          chan int
	totalTimes, finishTimes int
	simulateDayArr          []time.Time
	targetArrMap            targetArrMutex
	discardOverTime         bool
)

// Simulate Simulate
func Simulate() {
	prompt := promptui.Prompt{
		Label: "Balance type?(a: forward, b: reverse, c: force_both)",
	}
	balanceTypeAns, err := prompt.Run()
	if err != nil {
		panic(err)
	}
	switch balanceTypeAns {
	case "a":
		balanceType = "a"
	case "b":
		balanceType = "b"
	case "c":
		balanceType = "c"
	}

	prompt = promptui.Prompt{
		Label: "Discard over time trade?(y/n)",
	}
	result, err := prompt.Run()
	if err != nil {
		panic(err)
	}
	if result == "y" {
		discardOverTime = true
	}

	var useGlobal bool
	prompt = promptui.Prompt{
		Label: "Use global cond?(y/n)",
	}
	result, err = prompt.Run()
	if err != nil {
		panic(err)
	}
	if result == "y" {
		useGlobal = true
	}

	prompt = promptui.Prompt{
		Label: "N days?",
	}
	result, err = prompt.Run()
	if err != nil {
		panic(err)
	}
	n, err := common.StrToInt64(result)
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
					logger.GetLogger().Error(err)
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
	logger.GetLogger().Info("Fetch done")

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
	logger.GetLogger().Warn("Finish simulate")
	time.Sleep(10 * time.Second)
}

// getBestCond getBestCond
func getBestCond(historyCount int, useGlobal bool) {
	var wg sync.WaitGroup
	var conds []*simulationcond.AnalyzeCondition
	if useGlobal {
		conds = append(conds, &global.TickAnalyzeCondition)
	} else {
		for m := 95; m >= 80; m -= 5 {
			for u := 3; u <= 9; u += 3 {
				for i := 50; i >= 50; i -= 5 {
					for z := 0; z <= 5; z++ {
						for o := 10; o >= 10; o -= 2 {
							for p := 1; p <= 2; p++ {
								for v := 5; v <= 8; v++ {
									for g := -3; g <= -3; g++ {
										for h := 7; h >= 7; h-- {
											cond := simulationcond.AnalyzeCondition{
												HistoryCloseCount:    int64(historyCount),
												OutInRatio:           float64(m),
												ReverseOutInRatio:    float64(u),
												CloseDiff:            0,
												CloseChangeRatioLow:  float64(g),
												CloseChangeRatioHigh: float64(h),
												OpenChangeRatio:      float64(h),
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
		}
	}
	if err := simulationcond.InsertMultiRecord(conds, db.GetAgent()); err != nil {
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
				logger.GetLogger().Warnf("%s has no %s's close", stockNum, global.LastLastTradeDay.Format(global.ShortTimeLayout))
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
	var tradeCount, positiveCount int64
	for date, analyzeMap := range analyzeMapMap {
		var dateForwardBalance, dateReverseBalance int64
		if balanceType == "c" && (len(analyzeMap[0]) == 0 || len(analyzeMap[1]) == 0) && training {
			totalTimesChan <- -1
			return
		}
		sellTimeStamp := make(map[string]int64)
		for stockNum, v := range analyzeMap[0] {
			ticks := allTickMap.getAllTicksByStockNumAndDate(stockNum, date)
			endTradeInTime := getLastTradeInTimeByEntireTickTimeStamp(ticks[0].TimeStamp)
			endTradeOutTime := getLastTradeOutTimeByEntireTickTimeStamp(ticks[0].TimeStamp)
			var historyClose []float64
			var buyPrice, sellPrice float64
			for _, k := range ticks {
				historyClose = append(historyClose, k.Close)
				if len(historyClose) > int(cond.HistoryCloseCount) {
					historyClose = historyClose[1:]
				}
				if k.TimeStamp == v.TimeStamp && buyPrice == 0 && v.TimeStamp < endTradeInTime {
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
			if sellPrice == 0 && !training && buyPrice != 0 {
				logger.GetLogger().Warnf("%s no sell point", stockNum)
			} else {
				buyCost := tradebot.GetStockBuyCost(buyPrice, global.OneTimeQuantity)
				sellCost := tradebot.GetStockSellCost(sellPrice, global.OneTimeQuantity)
				back := tradebot.GetStockTradeFeeDiscount(buyPrice, global.OneTimeQuantity) + tradebot.GetStockTradeFeeDiscount(sellPrice, global.OneTimeQuantity)
				forwardBalance += (sellCost - buyCost + back)
				dateForwardBalance += (sellCost - buyCost + back)
				if sellTimeStamp[v.StockNum] > endTradeOutTime && training && discardOverTime {
					totalTimesChan <- -1
					return
				}
				tradeCount++
				if !training && (sellCost-buyCost+back) != 0 {
					buyTime := time.Unix(0, v.TimeStamp).Add(-8 * time.Hour)
					sellTime := time.Unix(0, sellTimeStamp[v.StockNum]).Add(-8 * time.Hour)
					logger.GetLogger().Warnf("%s Forward Balance: %d, Stock: %s, Name: %s, Buy at %s, Sell at %s", date, sellCost-buyCost+back, v.StockNum, global.AllStockNameMap.GetName(v.StockNum), buyTime.Format(global.LongTimeLayout), sellTime.Format(global.LongTimeLayout))
				}
			}
		}
		sellTimeStamp = make(map[string]int64)
		for stockNum, v := range analyzeMap[1] {
			ticks := allTickMap.getAllTicksByStockNumAndDate(stockNum, date)
			endTradeInTime := getLastTradeInTimeByEntireTickTimeStamp(ticks[0].TimeStamp)
			endTradeOutTime := getLastTradeOutTimeByEntireTickTimeStamp(ticks[0].TimeStamp)
			var historyClose []float64
			var sellFirstPrice, buyLaterPrice float64
			for _, k := range ticks {
				historyClose = append(historyClose, k.Close)
				if len(historyClose) > int(cond.HistoryCloseCount) {
					historyClose = historyClose[1:]
				}
				if k.TimeStamp == v.TimeStamp && sellFirstPrice == 0 && v.TimeStamp < endTradeInTime {
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
			if buyLaterPrice == 0 && !training && sellFirstPrice != 0 {
				logger.GetLogger().Warnf("%s no buy later point", stockNum)
			} else {
				buyCost := tradebot.GetStockBuyCost(buyLaterPrice, global.OneTimeQuantity)
				sellCost := tradebot.GetStockSellCost(sellFirstPrice, global.OneTimeQuantity)
				back := tradebot.GetStockTradeFeeDiscount(buyLaterPrice, global.OneTimeQuantity) + tradebot.GetStockTradeFeeDiscount(sellFirstPrice, global.OneTimeQuantity)
				if sellTimeStamp[v.StockNum] > endTradeOutTime && training && discardOverTime {
					totalTimesChan <- -1
					return
				}
				tradeCount++
				reverseBalance += (sellCost - buyCost + back)
				dateReverseBalance += (sellCost - buyCost + back)
				if !training && (sellCost-buyCost+back) != 0 {
					sellFirstTime := time.Unix(0, v.TimeStamp).Add(-8 * time.Hour)
					buyLaterTime := time.Unix(0, sellTimeStamp[v.StockNum]).Add(-8 * time.Hour)
					logger.GetLogger().Warnf("%s Reverse Balance: %d, Stock: %s, Name: %s, Sell first at %s, Buy later at %s", date, sellCost-buyCost+back, v.StockNum, global.AllStockNameMap.GetName(v.StockNum), sellFirstTime.Format(global.LongTimeLayout), buyLaterTime.Format(global.LongTimeLayout))
				}
			}
		}
		if !training {
			logger.GetLogger().Warnf("%s Forward: %d, Reverse: %d", date, dateForwardBalance, dateReverseBalance)
		}
		if dateForwardBalance+dateReverseBalance > 0 {
			positiveCount++
		}
	}
	tmp := simulate.Result{
		Balance:        forwardBalance + reverseBalance,
		ForwardBalance: forwardBalance,
		ReverseBalance: reverseBalance,
		TradeCount:     tradeCount,
		Cond:           cond,
		PositiveDays:   positiveCount,
		TotalDays:      int64(len(analyzeMapMap)),
	}
	if training {
		resultChan <- tmp
	} else {
		logger.GetLogger().Warnf("Total Balance: %d, TradeCount: %d, PositiveCount: %d", tmp.Balance, tradeCount, positiveCount)
		logger.GetLogger().Warnf("Cond: %+v", cond)
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
			if err := simulate.InsertMultiRecord(save, db.GetAgent()); err != nil {
				logger.GetLogger().Error(err)
			}
			var wg sync.WaitGroup
			if len(tmp) > 0 {
				logger.GetLogger().Info("Below is best result")
				for _, k := range tmp {
					wg.Add(1)
					go GetBalance(SearchTradePoint(simulateDayArr, k.Cond), k.Cond, false, &wg)
				}
				wg.Wait()
			} else if !useGlobal {
				logger.GetLogger().Info("No best result")
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
			// logger.GetLogger().Warnf("Best: %d", tmp[0].Balance)
			if err := simulate.InsertMultiRecord(save, db.GetAgent()); err != nil {
				logger.GetLogger().Error(err)
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
			// logger.GetLogger().Warnf("Total simulate counts: %d", totalTimes)
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
			ticks, err := entiretick.GetAllEntiretickByStockByDate(stockNum, date.Format(global.ShortTimeLayout), db.GetAgent())
			if err != nil {
				logger.GetLogger().Error(err)
				continue
			}
			allTickMap.saveByStockNumAndDate(stockNum, date.Format(global.ShortTimeLayout), ticks)
		}
	}
}

func clearAllSimulation() {
	if err := simulate.DeleteAll(db.GetAgent()); err != nil {
		panic(err)
	}
	if err := simulationcond.DeleteAll(db.GetAgent()); err != nil {
		panic(err)
	}
}
