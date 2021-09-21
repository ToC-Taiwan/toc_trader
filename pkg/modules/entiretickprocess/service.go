// Package entiretickprocess package entiretickprocess
package entiretickprocess

import (
	"errors"
	"runtime/debug"
	"sync"

	"github.com/markcheno/go-quote"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/analyzeentiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/entiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tickanalyze"
	"gitlab.tocraw.com/root/toc_trader/tools/common"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
)

// TickProcess TickProcess
func TickProcess(stockNum string, lastClose float64, cond global.AnalyzeCondition, ch chan *entiretick.EntireTick, wg *sync.WaitGroup, saveCh chan []*entiretick.EntireTick, sim bool, simulateMap *AnalyzeEntireTickMap) {
	var input quote.Quote
	var high, low, open float64
	var unSavedTicks entiretick.PtrArrArr
	var tmpArr entiretick.PtrArr
	var lastSaveLastClose float64
	var openChangeRatio float64
	analyzeChan := make(chan *analyzeentiretick.AnalyzeEntireTick)
	go AnalyzeEntireTickSaver(analyzeChan, wg, sim, simulateMap)
	for {
		tick, ok := <-ch
		if !ok {
			close(analyzeChan)
			break
		}
		if high == 0 && low == 0 && open == 0 {
			openChangeRatio = common.Round(100*(tick.Close-lastClose)/lastClose, 2)
			high = tick.Close
			low = tick.Close
			open = tick.Close
		}
		switch {
		case tick.Close <= high && tick.Close >= low:
			tick.High = high
			tick.Low = low
		case tick.Close > high:
			high = tick.Close
			tick.High = high
			tick.Low = low
		case tick.Close < low:
			low = tick.Close
			tick.Low = low
			tick.High = high
		}
		tick.Open = open
		tmpArr = append(tmpArr, tick)
		if tmpArr.GetTotalTime() < 20 {
			continue
		}
		if tmpArr.GetTotalTime() > 30 {
			unSavedTicks.ClearAll()
		}
		unSavedTicks.Append(tmpArr)
		if !sim {
			saveCh <- tmpArr
		}
		tmpArr = []*entiretick.EntireTick{}
		if unSavedTicks.GetCount() >= 3 {
			var outSum, inSum int64
			var totalTime float64
			closeChangeRatio := common.Round(100*(tick.Close-lastClose)/lastClose, 2)
			data := unSavedTicks.Get()
			for _, v := range data {
				input.Close = append(input.Close, v.GetAllCloseArr()...)
				outSum += v.GetOutSum()
				inSum += v.GetInSum()
				totalTime += v.GetTotalTime()
			}
			if len(input.Close) < global.HistoryCloseCount {
				unSavedTicks.ClearAll()
				continue
			} else {
				input.Close = input.Close[len(input.Close)-global.HistoryCloseCount:]
			}
			closeDiff := common.Round((unSavedTicks.GetLastClose() - lastSaveLastClose), 2)
			if lastSaveLastClose == 0 {
				closeDiff = 0
			}
			lastSaveLastClose = unSavedTicks.GetLastClose()
			unSavedTicksInOutRatio := common.Round(100*(float64(outSum)/float64(outSum+inSum)), 2)
			analyze := analyzeentiretick.AnalyzeEntireTick{
				TimeStamp:        tick.TimeStamp,
				StockNum:         stockNum,
				Close:            tick.Close,
				CloseChangeRatio: closeChangeRatio,
				OpenChangeRatio:  openChangeRatio,
				OutSum:           outSum,
				InSum:            inSum,
				OutInRatio:       unSavedTicksInOutRatio,
				TotalTime:        common.Round(totalTime, 2),
				CloseDiff:        closeDiff,
				Rsi:              tickanalyze.GenerateRSI(input),
				Open:             open,
				High:             high,
				Low:              low,
			}
			if unSavedTicksInOutRatio >= cond.OutInRatio && outSum >= cond.OutSum && closeDiff > cond.CloseDiff {
				if closeChangeRatio >= cond.CloseChangeRatioLow && closeChangeRatio <= cond.CloseChangeRatioHigh && closeChangeRatio <= cond.OpenChangeRatio && analyze.Rsi < float64(cond.RsiLow) {
					analyzeChan <- &analyze
					// name := global.AllStockNameMap.GetName(stockNum)
					// tickTime := time.Unix(0, tick.TimeStamp).UTC().Format(global.LongTimeLayout)
					// replaceDate := tickTime[:10]
					// clockTime := tickTime[11:19]
					// logger.Logger.WithFields(map[string]interface{}{
					// 	"Close":       analyze.Close,
					// 	"ChangeRatio": closeChangeRatio,
					// 	"OutSum":      outSum,
					// 	"InSum":       inSum,
					// 	"OutInRatio":  unSavedTicksInOutRatio,
					// 	"Name":        name,
					// 	"RSI":         analyze.Rsi,
					// }).Infof("EntireTick Analyze: %s %s %s", replaceDate, clockTime, stockNum)
				}
			}
			unSavedTicks.ClearAll()
		}
	}
}

// SaveEntireTicks SaveEntireTicks
func SaveEntireTicks(saveCh chan []*entiretick.EntireTick) {
	for {
		saveData, ok := <-saveCh
		if !ok {
			return
		}
		if len(saveData) != 0 {
			if err := entiretick.InsertMultiRecord(saveData, global.GlobalDB); err != nil {
				logger.Logger.Error(err)
				continue
			}
		}
	}
}

// AnalyzeEntireTickSaver AnalyzeEntireTickSaver
func AnalyzeEntireTickSaver(ch chan *analyzeentiretick.AnalyzeEntireTick, wg *sync.WaitGroup, sim bool, simulateMap *AnalyzeEntireTickMap) {
	var err error
	defer wg.Done()
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
			logger.Logger.Error(err.Error() + "\n" + string(debug.Stack()))
		}
	}()
	var tmpArr []*analyzeentiretick.AnalyzeEntireTick
	for {
		tick, ok := <-ch
		if !ok {
			if !sim {
				if len(tmpArr) != 0 {
					if err := analyzeentiretick.InsertMultiRecord(tmpArr, global.GlobalDB); err != nil {
						panic(err)
					}
				}
			} else {
				if len(tmpArr) != 0 {
					simulateMap.SaveByStockNum(tmpArr[0].StockNum, tmpArr)
				}
			}
			break
		}
		tmpArr = append(tmpArr, tick)
	}
}
