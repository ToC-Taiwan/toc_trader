// Package entiretickprocess package entiretickprocess
package entiretickprocess

import (
	"errors"
	"runtime/debug"
	"sync"

	"github.com/markcheno/go-quote"
	"gitlab.tocraw.com/root/toc_trader/internal/common"
	"gitlab.tocraw.com/root/toc_trader/internal/database"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/analyzeentiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/entiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tickanalyze"
)

// TickProcess TickProcess
func TickProcess(stockNum string, lastClose float64, cond simulationcond.AnalyzeCondition, ch chan *entiretick.EntireTick, wg *sync.WaitGroup, saveCh chan []*entiretick.EntireTick, sim bool, simulateMap *AnalyzeEntireTickMap) {
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
		if tmpArr.GetTotalTime() < cond.TicksPeriodThreshold {
			continue
		}
		if tmpArr.GetTotalTime() > cond.TicksPeriodLimit {
			unSavedTicks.ClearAll()
		}
		unSavedTicks.Append(tmpArr)
		if !sim {
			saveCh <- tmpArr
		}
		tmpArr = []*entiretick.EntireTick{}
		if unSavedTicks.GetCount() >= cond.TicksPeriodCount {
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
			if len(input.Close) < int(cond.HistoryCloseCount) {
				unSavedTicks.ClearAll()
				continue
			} else if cond.TrimHistoryCloseCount {
				input.Close = input.Close[len(input.Close)-int(cond.HistoryCloseCount):]
			}
			closeDiff := common.Round((unSavedTicks.GetLastClose() - lastSaveLastClose), 2)
			if lastSaveLastClose == 0 {
				closeDiff = 0
			}
			lastSaveLastClose = unSavedTicks.GetLastClose()
			unSavedTicksInOutRatio := common.Round(100*(float64(outSum)/float64(outSum+inSum)), 2)
			rsi, err := tickanalyze.GenerateRSI(input)
			if err != nil {
				logger.GetLogger().Errorf("GenerateRSI at EntireTickProcess Stock: %s, Err: %s", stockNum, err)
				continue
			}
			analyze := analyzeentiretick.AnalyzeEntireTick{
				TimeStamp:        tick.TimeStamp,
				StockNum:         stockNum,
				Close:            tick.Close,
				CloseChangeRatio: closeChangeRatio,
				OpenChangeRatio:  openChangeRatio,
				OutSum:           outSum,
				InSum:            inSum,
				OutInRatio:       unSavedTicksInOutRatio,
				TotalTime:        totalTime,
				CloseDiff:        closeDiff,
				Rsi:              rsi,
				Open:             open,
				High:             high,
				Low:              low,
				Volume:           outSum + inSum,
			}
			analyzeChan <- &analyze
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
			if err := entiretick.InsertMultiRecord(saveData, database.GetAgent()); err != nil {
				logger.GetLogger().Error(err)
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
			logger.GetLogger().Error(err.Error() + "\n" + string(debug.Stack()))
		}
	}()
	var tmpArr []*analyzeentiretick.AnalyzeEntireTick
	for {
		tick, ok := <-ch
		if !ok {
			if !sim {
				if len(tmpArr) != 0 {
					if err := analyzeentiretick.InsertMultiRecord(tmpArr, database.GetAgent()); err != nil {
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
