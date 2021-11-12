// Package simulateprocess package simulateprocess
package simulateprocess

import (
	"sync"
	"time"

	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
	"gitlab.tocraw.com/root/toc_trader/internal/common"
	"gitlab.tocraw.com/root/toc_trader/internal/database"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/analyzeentiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/entiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulate"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/choosetarget"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/fetchentiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/importbasic"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tickprocess"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tradebot"
)

var (
	balanceType             simulateType
	allTickMap              entireTickMap
	resultChan              chan simulate.Result
	totalTimesChan          chan int
	totalTimes, finishTimes int
	simulateDayArr          []time.Time
	targetArrMap            targetArrMutex
	discardOverTime         bool
)

// Simulate Simulate
func Simulate(simType, discardOT, useDefault, dayCount string) {
	allTickMap.clearAll()
	targetArrMap.clearAll()
	totalTimes = 0
	finishTimes = 0
	switch simType {
	case "a":
		balanceType = simTypeForward
	case "b":
		balanceType = simTypeReverse
	}
	if discardOT == "y" {
		discardOverTime = true
	}
	var useGlobal bool
	if useDefault == "y" {
		useGlobal = true
	}
	n, err := common.StrToInt64(dayCount)
	if err != nil {
		panic(err)
	}
	// clearAllSimulation()
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
				logger.GetLogger().WithFields(map[string]interface{}{
					"Date": date.Format(global.ShortTimeLayout),
					"Rank": i + 1,
					"Name": global.AllStockNameMap.GetName(v),
				}).Infof("Volume Rank")
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
			fetchentiretick.FetchEntireTick(targets, tmp, global.BaseCond)
			storeAllEntireTick(targets, tmp)
		}
	}
	simulateDayArr = tradeDayArr
	logger.GetLogger().Info("Fetch Done")

	resultChan = make(chan simulate.Result)
	totalTimesChan = make(chan int)
	go totalTimesReceiver()
	go catchResult(useGlobal)

	var wg sync.WaitGroup
	if useGlobal {
		var historyCount int64
		switch balanceType {
		case simTypeForward:
			historyCount = global.ForwardCond.HistoryCloseCount
		case simTypeReverse:
			historyCount = global.ForwardCond.HistoryCloseCount
		}
		getBestCond(int(historyCount), useGlobal)
	} else {
		for i := 2500; i >= 100; i -= 100 {
			wg.Add(1)
			go func(historyCount int) {
				defer wg.Done()
				getBestCond(historyCount, useGlobal)
			}(i)
		}
	}
	wg.Wait()
	close(resultChan)
	logger.GetLogger().Info("Finish Simulate")
	time.Sleep(10 * time.Second)
}

// getBestCond getBestCond
func getBestCond(historyCount int, useGlobal bool) {
	var wg sync.WaitGroup
	var conds []*simulationcond.AnalyzeCondition
	if useGlobal {
		switch balanceType {
		case simTypeForward:
			conds = append(conds, &global.ForwardCond)
		case simTypeReverse:
			conds = append(conds, &global.ReverseCond)
		}
	} else {
		switch balanceType {
		case simTypeForward:
			conds = generateForwardConds(historyCount)
		case simTypeReverse:
			conds = generateReverseConds(historyCount)
		}
	}
	if err := simulationcond.InsertMultiRecord(conds, database.GetAgent()); err != nil {
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
	var simulateAnalyzeEntireMap tickprocess.AnalyzeEntireTickMap
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
				go tickprocess.TickProcess(stockNum, lastClose, cond, ch, &wg, saveCh, true, &simulateAnalyzeEntireMap)
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
			lastTime := time.Date(tickTimeUnix.Year(), tickTimeUnix.Month(), tickTimeUnix.Day(), global.TradeOutEndHour, global.TradeOutEndMinute, 0, 0, time.Local)
			if tickTimeUnix.After(lastTime) || buyPointMap[v.StockNum] != nil || sellFirstPointMap[v.StockNum] != nil {
				continue
			}
			if tradebot.IsBuyPoint(tmp, cond) && balanceType == simTypeForward {
				buyPointMap[v.StockNum] = v
				continue
			}
			if tradebot.IsSellFirstPoint(tmp, cond) && balanceType == simTypeReverse {
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
	var forwardBalance, reverseBalance, totalLoss int64
	var tradeCount, positiveCount int64
	for date, analyzeMap := range analyzeMapMap {
		sellTimeStamp := make(map[string]int64)
		var dateForwardBalance, dateReverseBalance int64
		for stockNum, v := range analyzeMap[0] {
			ticks := allTickMap.getAllTicksByStockNumAndDate(stockNum, date)
			endTradeInTime := getLastTradeInTimeByEntireTickTimeStamp(ticks[0].TimeStamp)
			endTradeOutTime := getLastTradeOutTimeByEntireTickTimeStamp(ticks[0].TimeStamp)
			var historyClose []float64
			var buyPrice, sellPrice float64
			for _, k := range ticks {
				historyClose = append(historyClose, k.Close)
				if len(historyClose) > int(cond.HistoryCloseCount) && cond.TrimHistoryCloseCount {
					historyClose = historyClose[1:]
				}
				if k.TimeStamp == v.TimeStamp && buyPrice == 0 && v.TimeStamp < endTradeInTime {
					buyPrice = k.Close
					tradeCount++
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
				tmpBalance := (sellCost - buyCost + back)
				forwardBalance += tmpBalance
				if tmpBalance < 0 {
					totalLoss += tmpBalance
				}
				dateForwardBalance += tmpBalance
				if sellTimeStamp[v.StockNum] > endTradeOutTime && training && discardOverTime {
					totalTimesChan <- -1
					return
				}
				if !training && (sellCost-buyCost+back) != 0 {
					buyTime := time.Unix(0, v.TimeStamp).Add(-8 * time.Hour)
					sellTime := time.Unix(0, sellTimeStamp[v.StockNum]).Add(-8 * time.Hour)
					logger.GetLogger().WithFields(map[string]interface{}{
						"Date":            date,
						"OriginalBalance": sellCost - buyCost,
						"Back":            back,
						"Name":            global.AllStockNameMap.GetName(v.StockNum),
						"BuyAt":           buyTime.Format(global.LongTimeLayout)[11:],
						"SellAt":          sellTime.Format(global.LongTimeLayout)[11:],
					}).Warn("Forward Balance")
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
				if len(historyClose) > int(cond.HistoryCloseCount) && cond.TrimHistoryCloseCount {
					historyClose = historyClose[1:]
				}
				if k.TimeStamp == v.TimeStamp && sellFirstPrice == 0 && v.TimeStamp < endTradeInTime {
					sellFirstPrice = k.Close
					tradeCount++
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
				tmpBalance := (sellCost - buyCost + back)
				if tmpBalance < 0 {
					totalLoss += tmpBalance
				}
				reverseBalance += tmpBalance
				dateReverseBalance += tmpBalance
				if !training && (sellCost-buyCost+back) != 0 {
					sellFirstTime := time.Unix(0, v.TimeStamp).Add(-8 * time.Hour)
					buyLaterTime := time.Unix(0, sellTimeStamp[v.StockNum]).Add(-8 * time.Hour)
					logger.GetLogger().WithFields(map[string]interface{}{
						"Date":            date,
						"OriginalBalance": sellCost - buyCost,
						"Back":            back,
						"Name":            global.AllStockNameMap.GetName(v.StockNum),
						"SellFirstAt":     sellFirstTime.Format(global.LongTimeLayout)[11:],
						"BuyLaterAt":      buyLaterTime.Format(global.LongTimeLayout)[11:],
					}).Warn("Reverse Balance")
				}
			}
		}
		if !training {
			logger.GetLogger().WithFields(map[string]interface{}{
				"ForwardBalance": dateForwardBalance,
				"ReverseBalance": dateReverseBalance,
				"Date":           date,
			}).Warn("Single Day Summary")
		}
		if dateForwardBalance+dateReverseBalance > 0 {
			positiveCount++
		}
	}
	tmp := simulate.Result{
		Balance:        forwardBalance + reverseBalance,
		ForwardBalance: forwardBalance,
		ReverseBalance: reverseBalance,
		TotalLoss:      totalLoss * -1,
		TradeCount:     tradeCount,
		Cond:           cond,
		PositiveDays:   positiveCount,
		TotalDays:      int64(len(analyzeMapMap)),
	}
	if training {
		resultChan <- tmp
	} else if tmp.Balance != 0 {
		logger.GetLogger().WithFields(map[string]interface{}{
			"TradeCount":    tradeCount,
			"Balance":       tmp.Balance,
			"PositiveCount": positiveCount,
		}).Warn("Total Balance")
		logger.GetLogger().WithFields(map[string]interface{}{
			"TrimHistoryCloseCount": cond.TrimHistoryCloseCount,
			"HistoryCloseCount":     cond.HistoryCloseCount,
		}).Warn("Cond")
		logger.GetLogger().WithFields(map[string]interface{}{
			"ForwardOutInRatio": cond.ForwardOutInRatio,
			"ReverseOutInRatio": cond.ReverseOutInRatio,
		}).Warn("Cond")
		logger.GetLogger().WithFields(map[string]interface{}{
			"CloseChangeRatioLow":  cond.CloseChangeRatioLow,
			"CloseChangeRatioHigh": cond.CloseChangeRatioHigh,
			"OpenChangeRatio":      cond.OpenChangeRatio,
		}).Warn("Cond")
		logger.GetLogger().WithFields(map[string]interface{}{
			"RsiHigh": cond.RsiHigh,
			"RsiLow":  cond.RsiLow,
		}).Warn("Cond")
		logger.GetLogger().WithFields(map[string]interface{}{
			"TicksPeriodThreshold": cond.TicksPeriodThreshold,
			"TicksPeriodLimit":     cond.TicksPeriodLimit,
			"TicksPeriodCount":     cond.TicksPeriodCount,
			"VolumePerSecond":      cond.VolumePerSecond,
		}).Warn("Cond")
	}
}

func catchResult(useGlobal bool) {
	var save []simulate.Result
	var count int
	for {
		result, ok := <-resultChan
		if result.Cond.Model.ID != 0 && result.TradeCount != 0 {
			save = append(save, result)
		}
		if !ok {
			if err := simulate.InsertMultiRecord(save, database.GetAgent()); err != nil {
				logger.GetLogger().Error(err)
			}
			bestResult, err := simulate.GetBestSimulateResult(database.GetAgent())
			if err != nil {
				panic(err)
			}
			if balanceType == simTypeForward {
				bestResult.IsBestForward = true
			} else if balanceType == simTypeReverse {
				bestResult.IsBestReverse = true
			}
			if err := simulate.Update(&bestResult, database.GetAgent()); err != nil {
				panic(err)
			}
			if bestResult.ID != 0 {
				var wg sync.WaitGroup
				logger.GetLogger().Info("Best Result")
				wg.Add(1)
				go GetBalance(SearchTradePoint(simulateDayArr, bestResult.Cond), bestResult.Cond, false, &wg)
				wg.Wait()
			} else if !useGlobal {
				logger.GetLogger().Info("No Best")
			}
			close(totalTimesChan)
			break
		}
		count++
		totalTimesChan <- -1
		if count%10 == 0 {
			if err := simulate.InsertMultiRecord(save, database.GetAgent()); err != nil {
				logger.GetLogger().Error(err)
			}
			save = []simulate.Result{}
		}
	}
}

func totalTimesReceiver() {
	var count int
	for {
		times, ok := <-totalTimesChan
		if !ok {
			break
		}
		if times > 0 {
			count++
			totalTimes += times
			if count == 25 {
				go progressBar(totalTimes)
			}
		} else {
			finishTimes++
			if count == 25 {
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
		mpb.AppendDecorators(decor.Counters(0, "")),
	)
	p.Wait()
}

func storeAllEntireTick(stockArr []string, tradeDayArr []time.Time) {
	for _, stockNum := range stockArr {
		for _, date := range tradeDayArr {
			ticks, err := entiretick.GetAllEntiretickByStockByDate(stockNum, date.Format(global.ShortTimeLayout), database.GetAgent())
			if err != nil {
				logger.GetLogger().Error(err)
				continue
			}
			allTickMap.saveByStockNumAndDate(stockNum, date.Format(global.ShortTimeLayout), ticks)
		}
	}
}

// ClearAllSimulation ClearAllSimulation
func ClearAllSimulation() {
	if err := simulate.DeleteAll(database.GetAgent()); err != nil {
		panic(err)
	}
	if err := simulationcond.DeleteAll(database.GetAgent()); err != nil {
		panic(err)
	}
}

func generateForwardConds(historyCount int) []*simulationcond.AnalyzeCondition {
	var conds []*simulationcond.AnalyzeCondition
	var i float64
	for m := 95; m >= 85; m -= 5 {
		for g := -1; g <= -1; g++ {
			for h := 3; h >= 3; h-- {
				for i = 0.9; i >= 0.6; i -= 0.1 {
					for o := 12; o >= 4; o -= 4 {
						for p := 3; p >= 1; p-- {
							for v := 30; v >= 10; v -= 10 {
								for r := 1; r <= 3; r++ {
									cond := simulationcond.AnalyzeCondition{
										TrimHistoryCloseCount: true,
										HistoryCloseCount:     int64(historyCount),
										ForwardOutInRatio:     float64(m),
										CloseChangeRatioLow:   float64(g),
										CloseChangeRatioHigh:  float64(h),
										OpenChangeRatio:       float64(g),
										RsiHigh:               i,
										TicksPeriodThreshold:  float64(o),
										TicksPeriodLimit:      float64(o) * 1.3,
										TicksPeriodCount:      p,
										VolumePerSecond:       int64(v),
										MaxHoldTime:           int64(r),
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
	return conds
}

func generateReverseConds(historyCount int) []*simulationcond.AnalyzeCondition {
	var conds []*simulationcond.AnalyzeCondition
	var k float64
	for u := 3; u <= 9; u += 3 {
		for g := 0; g <= 0; g++ {
			for h := 3; h >= 3; h-- {
				for k = 0.1; k <= 0.4; k += 0.1 {
					for o := 12; o >= 4; o -= 4 {
						for p := 3; p >= 1; p-- {
							for v := 30; v >= 10; v -= 10 {
								for r := 1; r <= 3; r++ {
									cond := simulationcond.AnalyzeCondition{
										TrimHistoryCloseCount: true,
										HistoryCloseCount:     int64(historyCount),
										ReverseOutInRatio:     float64(u),
										CloseChangeRatioLow:   float64(g),
										CloseChangeRatioHigh:  float64(h),
										OpenChangeRatio:       float64(h),
										RsiLow:                k,
										TicksPeriodThreshold:  float64(o),
										TicksPeriodLimit:      float64(o) * 1.3,
										TicksPeriodCount:      p,
										VolumePerSecond:       int64(v),
										MaxHoldTime:           int64(r),
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
	return conds
}
