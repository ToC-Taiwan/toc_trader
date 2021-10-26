// Package globalinit is init all global var
package globalinit

import (
	"os"
	"time"

	"gitlab.tocraw.com/root/toc_trader/init/sysparminit"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/importbasic"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
)

func init() {
	var err error
	global.ExitChannel = make(chan int)
	global.HTTPPort = sysparminit.GlobalSettings.GetHTTPPort()
	global.PyServerHost = sysparminit.GlobalSettings.GetPyServerHost()
	global.PyServerPort = sysparminit.GlobalSettings.GetPyServerPort()

	global.TradeSwitch = global.SystemSwitch{
		Buy:                          true,
		Sell:                         true,
		SellFirst:                    true,
		BuyLater:                     true,
		UseBidAsk:                    false,
		MeanTimeTradeStockNum:        25,
		MeanTimeReverseTradeStockNum: 25,
	}

	deployment := os.Getenv("DEPLOYMENT")
	if deployment != "docker" {
		global.TradeSwitch.Buy = false
		global.TradeSwitch.SellFirst = false
	}

	global.CentralCond = simulationcond.AnalyzeCondition{
		TrimHistoryCloseCount: true,
		HistoryCloseCount:     2500,
		OutInRatio:            60,
		ReverseOutInRatio:     10,
		CloseDiff:             0,
		CloseChangeRatioLow:   -3,
		CloseChangeRatioHigh:  7,
		OpenChangeRatio:       7,
		RsiHigh:               50,
		RsiLow:                50,
		ReverseRsiHigh:        50,
		ReverseRsiLow:         50,
		TicksPeriodThreshold:  10,
		TicksPeriodLimit:      10 * 1.3,
		TicksPeriodCount:      1,
		VolumePerSecond:       5,
	}

	global.ForwardCond = simulationcond.AnalyzeCondition{
		TrimHistoryCloseCount: true,
		HistoryCloseCount:     2500,
		OutInRatio:            90,
		ReverseOutInRatio:     3,
		CloseDiff:             0,
		CloseChangeRatioLow:   -3,
		CloseChangeRatioHigh:  7,
		OpenChangeRatio:       7,
		RsiHigh:               50.4,
		RsiLow:                50,
		ReverseRsiHigh:        50.4,
		ReverseRsiLow:         50,
		TicksPeriodThreshold:  10,
		TicksPeriodLimit:      10 * 1.3,
		TicksPeriodCount:      2,
		VolumePerSecond:       5,
	}

	global.ReverseCond = simulationcond.AnalyzeCondition{
		TrimHistoryCloseCount: true,
		HistoryCloseCount:     2500,
		OutInRatio:            90,
		ReverseOutInRatio:     6,
		CloseDiff:             0,
		CloseChangeRatioLow:   -3,
		CloseChangeRatioHigh:  7,
		OpenChangeRatio:       7,
		RsiHigh:               50.4,
		RsiLow:                50,
		ReverseRsiHigh:        50.4,
		ReverseRsiLow:         50,
		TicksPeriodThreshold:  10,
		TicksPeriodLimit:      10 * 1.3,
		TicksPeriodCount:      1,
		VolumePerSecond:       10,
	}

	if err = importbasic.ImportHoliday(); err != nil {
		panic(err)
	}
	global.TradeDay, err = importbasic.GetTradeDay()
	if err != nil {
		panic(err)
	}
	global.TradeDayInEndTime = time.Date(
		global.TradeDay.Year(),
		global.TradeDay.Month(),
		global.TradeDay.Day(),
		global.TradeInEndHour,
		global.TradeInEndMinute,
		0,
		0,
		time.Local)
	global.TradeDayOutEndTime = time.Date(
		global.TradeDay.Year(),
		global.TradeDay.Month(),
		global.TradeDay.Day(),
		global.TradeOutEndHour,
		global.TradeOutEndMinute,
		0,
		0,
		time.Local)
	lastTradeDayArr, err := importbasic.GetLastNTradeDay(2)
	if err != nil {
		panic(err)
	}
	global.LastTradeDay = lastTradeDayArr[0]
	global.LastLastTradeDay = lastTradeDayArr[1]
	global.LastTradeDayArr = append(global.LastTradeDayArr, global.LastTradeDay, global.LastLastTradeDay)
	logger.GetLogger().WithFields(map[string]interface{}{
		"TradeDay":         global.TradeDay.Format(global.ShortTimeLayout),
		"LastTradeDay":     global.LastTradeDay.Format(global.ShortTimeLayout),
		"LastLastTradeDay": global.LastLastTradeDay.Format(global.ShortTimeLayout),
	}).Info("Last Trade Days")

	logger.GetLogger().WithFields(map[string]interface{}{
		"TradeDayInEndTime": global.TradeDayInEndTime.Format(global.LongTimeLayout),
	}).Info("Trade In End Time")
	logger.GetLogger().WithFields(map[string]interface{}{
		"TradeDayInEndTime": global.TradeDayOutEndTime.Format(global.LongTimeLayout),
	}).Info("Trade Out End Time")
}
