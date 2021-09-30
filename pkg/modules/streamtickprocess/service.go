// Package streamtickprocess package streamtickprocess
package streamtickprocess

import (
	"github.com/markcheno/go-quote"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/analyzestreamtick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/streamtick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tickanalyze"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tradebot"
	"gitlab.tocraw.com/root/toc_trader/tools/common"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
)

// TickProcess TickProcess
func TickProcess(lastClose float64, cond global.AnalyzeCondition, ch chan *streamtick.StreamTick, saveCh chan []*streamtick.StreamTick) {
	if lastClose == 0 {
		return
	}
	buyChan := make(chan *analyzestreamtick.AnalyzeStreamTick)
	sellChan := make(chan *streamtick.StreamTick)
	// analyzeChan := make(chan *analyzestreamtick.AnalyzeStreamTick)
	go tradebot.BuyBot(buyChan)
	go tradebot.SellBot(sellChan)
	// go AnalyzeStreamTickSaver(analyzeChan)

	var input quote.Quote
	var unSavedTicks streamtick.PtrArrArr
	var tmpArr streamtick.PtrArr
	var lastSaveLastClose float64
	for {
		tick := <-ch
		if tradebot.BuyOrderMap.CheckStockExist(tick.StockNum) {
			sellChan <- tick
		}
		tmpArr = append(tmpArr, tick)
		if tmpArr.GetTotalTime() < cond.TicksPeriodThreshold {
			continue
		}
		if tmpArr.GetTotalTime() > cond.TicksPeriodLimit {
			unSavedTicks.ClearAll()
		}
		unSavedTicks.Append(tmpArr)
		saveCh <- tmpArr
		tmpArr = []*streamtick.StreamTick{}

		if unSavedTicks.GetCount() >= cond.TicksPeriodCount {
			var outSum, inSum int64
			var totalTime float64
			data := unSavedTicks.Get()
			for _, v := range data {
				input.Close = append(input.Close, v.GetAllCloseArr()...)
				outSum += v.GetOutSum()
				inSum += v.GetInSum()
				totalTime += v.GetTotalTime()
			}
			if len(input.Close) < int(global.TickAnalyzeCondition.HistoryCloseCount) {
				unSavedTicks.ClearAll()
				continue
			} else {
				input.Close = input.Close[len(input.Close)-int(global.TickAnalyzeCondition.HistoryCloseCount):]
			}
			closeDiff := common.Round((unSavedTicks.GetLastClose() - lastSaveLastClose), 2)
			if lastSaveLastClose == 0 {
				closeDiff = 0
			}
			lastSaveLastClose = unSavedTicks.GetLastClose()
			unSavedTicksInOutRatio := common.Round((100 * (float64(outSum) / float64(outSum+inSum))), 2)
			rsi, err := tickanalyze.GenerateRSI(input)
			if err != nil {
				logger.Logger.Errorf("TickProcess Stock: %s, Err: %s", tick.StockNum, err)
				continue
			}
			analyze := analyzestreamtick.AnalyzeStreamTick{
				TimeStamp:        tick.TimeStamp,
				StockNum:         tick.StockNum,
				Close:            tick.Close,
				OpenChangeRatio:  common.Round((tick.Open - lastClose), 2),
				CloseChangeRatio: tick.PctChg,
				OutSum:           outSum,
				InSum:            inSum,
				OutInRatio:       unSavedTicksInOutRatio,
				TotalTime:        common.Round(totalTime, 2),
				CloseDiff:        closeDiff,
				Open:             tick.Open,
				AvgPrice:         tick.AvgPrice,
				High:             tick.High,
				Low:              tick.Low,
				Rsi:              rsi,
			}
			// analyzeChan <- analyze
			buyChan <- &analyze
			unSavedTicks.ClearAll()
		}
	}
}

// SaveStreamTicks SaveStreamTicks
func SaveStreamTicks(saveCh chan []*streamtick.StreamTick) {
	for {
		unSavedTicks := <-saveCh
		if len(unSavedTicks) != 0 {
			if err := streamtick.InsertMultiRecord(unSavedTicks, global.GlobalDB); err != nil {
				logger.Logger.Error(err)
				continue
			}
		}
	}
}

// var analyzeStreamTickTmpArr analyzestreamtick.AnalyzeStreamArrMutexStruct

// // AnalyzeStreamTickSaver AnalyzeStreamTickSaver
// func AnalyzeStreamTickSaver(ch chan *analyzestreamtick.AnalyzeStreamTick) {
// 	var err error
// 	defer func() {
// 		if r := recover(); r != nil {
// 			switch x := r.(type) {
// 			case string:
// 				err = errors.New(x)
// 			case error:
// 				err = x
// 			default:
// 				err = errors.New("unknown panic")
// 			}
// 			logger.Logger.Error(err.Error() + "\n" + string(debug.Stack()))
// 		}
// 	}()
// 	go func() {
// 		tick := time.Tick(5 * time.Second)
// 		for range tick {
// 			if analyzeStreamTickTmpArr.GetTotalCount() != 0 {
// 				analyzeStreamTickTmpArr.Mutex.Lock()
// 				if err := analyzestreamtick.InsertMultiRecord(analyzeStreamTickTmpArr.Ticks, global.GlobalDB); err != nil {
// 					panic(err)
// 				}
// 				analyzeStreamTickTmpArr.Ticks = []*analyzestreamtick.AnalyzeStreamTick{}
// 				analyzeStreamTickTmpArr.Mutex.Unlock()
// 			}
// 		}
// 	}()
// 	for {
// 		tick := <-ch
// 		analyzeStreamTickTmpArr.Append(tick)
// 	}
// }
