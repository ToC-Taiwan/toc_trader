// Package streamtickprocess package streamtickprocess
package streamtickprocess

import (
	"github.com/markcheno/go-quote"
	"gitlab.tocraw.com/root/toc_trader/internal/common"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/analyzestreamtick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/streamtick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tickanalyze"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tradebot"
)

// ReverseTickProcess ReverseTickProcess
func ReverseTickProcess(lastClose float64, cond simulationcond.AnalyzeCondition, ch chan *streamtick.StreamTick) {
	var input quote.Quote
	var unSavedTicks streamtick.PtrArrArr
	var tmpArr streamtick.PtrArr
	var lastSaveLastClose, openChangeRatio float64
	var buyLaterChan chan *streamtick.StreamTick
	if lastClose == 0 {
		return
	}
	analyzeTickChan := make(chan *analyzestreamtick.AnalyzeStreamTick)
	go tradebot.SellFirstAgent(analyzeTickChan)

	if global.TradeSwitch.BuyLater {
		buyLaterChan = make(chan *streamtick.StreamTick)
		go tradebot.BuyLaterBot(buyLaterChan, cond, &input.Close)
	}
	for {
		tick := <-ch
		if openChangeRatio == 0 {
			openChangeRatio = common.Round((tick.Open - lastClose), 2)
		}
		tmpArr = append(tmpArr, tick)
		if tradebot.FilledSellFirstOrderMap.CheckStockExist(tick.StockNum) {
			buyLaterChan <- tick
		}

		if tmpArr.GetTotalTime() < cond.TicksPeriodThreshold {
			continue
		}
		if tmpArr.GetTotalTime() > cond.TicksPeriodLimit {
			unSavedTicks.ClearAll()
		}
		unSavedTicks.Append(tmpArr)
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
			if len(input.Close) < int(cond.HistoryCloseCount) {
				unSavedTicks.ClearAll()
				continue
			} else if cond.TrimHistoryCloseCount {
				input.Close = input.Close[len(input.Close)-int(cond.HistoryCloseCount):]
			}
			rsi, err := tickanalyze.GenerateRSI(input)
			if err != nil {
				logger.GetLogger().Errorf("GenerateRSI at StreamTickProcess Stock: %s, Err: %s", tick.StockNum, err)
				continue
			}

			closeDiff := common.Round((unSavedTicks.GetLastClose() - lastSaveLastClose), 2)
			if lastSaveLastClose == 0 {
				closeDiff = 0
			}
			lastSaveLastClose = unSavedTicks.GetLastClose()
			unSavedTicksInOutRatio := common.Round((100 * (float64(outSum) / float64(outSum+inSum))), 2)
			analyze := analyzestreamtick.AnalyzeStreamTick{
				TimeStamp:        tick.TimeStamp,
				StockNum:         tick.StockNum,
				Close:            tick.Close,
				OpenChangeRatio:  openChangeRatio,
				CloseChangeRatio: tick.PctChg,
				OutSum:           outSum,
				InSum:            inSum,
				OutInRatio:       unSavedTicksInOutRatio,
				TotalTime:        totalTime,
				CloseDiff:        closeDiff,
				Open:             tick.Open,
				AvgPrice:         tick.AvgPrice,
				High:             tick.High,
				Low:              tick.Low,
				Rsi:              rsi,
				Volume:           outSum + inSum,
			}
			analyzeTickChan <- &analyze
			unSavedTicks.ClearAll()
		}
	}
}
