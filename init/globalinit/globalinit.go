// Package globalinit is init all global var
package globalinit

import (
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

	global.TickAnalyzeCondition = simulationcond.AnalyzeCondition{
		HistoryCloseCount:    2500,
		OutInRatio:           80,
		ReverseOutInRatio:    5,
		CloseDiff:            0,
		CloseChangeRatioLow:  0,
		CloseChangeRatioHigh: 6,
		OpenChangeRatio:      3,
		RsiHigh:              50.1,
		RsiLow:               50,
		ReverseRsiHigh:       50.1,
		ReverseRsiLow:        50,
		TicksPeriodThreshold: 5,
		TicksPeriodLimit:     5 * 1.3,
		TicksPeriodCount:     1,
		VolumePerSecond:      5,
	}

	if err = importbasic.ImportHoliday(); err != nil {
		panic(err)
	}
	global.TradeDay, err = importbasic.GetTradeDay()
	if err != nil {
		panic(err)
	}
	global.TradeInDayEndTime = time.Date(
		global.TradeDay.Year(),
		global.TradeDay.Month(),
		global.TradeDay.Day(),
		global.TradeInEndHour,
		global.TradeInEndMinute,
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
		"TradeInDayEndTime": global.TradeInDayEndTime.Format(global.LongTimeLayout),
	}).Info("Trade End Time")
}
