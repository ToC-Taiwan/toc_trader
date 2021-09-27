// Package globalinit is init all global var
package globalinit

import (
	"time"

	"gitlab.tocraw.com/root/toc_trader/init/sysparminit"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/importbasic"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
)

func init() {
	global.ExitChannel = make(chan string)

	global.HTTPPort = sysparminit.GlobalSettings.GetHTTPPort()
	global.PyServerHost = sysparminit.GlobalSettings.GetPyServerHost()
	global.PyServerPort = sysparminit.GlobalSettings.GetPyServerPort()

	global.EnableBuy = true
	global.EnableSell = true
	global.UseBidAsk = false
	global.MeanTimeTradeStockNum = 3

	global.TickAnalyzeCondition = global.AnalyzeCondition{
		HistoryCloseCount:    500,
		OutSum:               100,
		OutInRatio:           70,
		CloseDiff:            0,
		CloseChangeRatioLow:  -1,
		CloseChangeRatioHigh: 10,
		OpenChangeRatio:      10,
		RsiHigh:              55,
		RsiLow:               55,
	}

	if err := importbasic.ImportHoliday(); err != nil {
		panic(err)
	}

	var today time.Time
	if time.Now().Hour() >= 15 {
		today = time.Now().AddDate(0, 0, 1)
	} else {
		today = time.Now()
	}
	tradeDay, err := importbasic.GetTradeDayTime(today)
	if err != nil {
		panic(err)
	}
	lastTradeDay, err := importbasic.GetLastTradeDayTime(tradeDay)
	if err != nil {
		panic(err)
	}

	lastLastTradeDay, err := importbasic.GetLastTradeDayTime(lastTradeDay)
	if err != nil {
		panic(err)
	}

	global.TradeDay = tradeDay
	global.LastTradeDay = lastTradeDay
	global.LastLastTradeDay = lastLastTradeDay
	global.LastTradeDayArr = append(global.LastTradeDayArr, lastLastTradeDay, lastTradeDay)

	logger.Logger.WithFields(map[string]interface{}{
		"TradeDay":         tradeDay.Format(global.ShortTimeLayout),
		"LastTradeDay":     lastTradeDay.Format(global.ShortTimeLayout),
		"LastLastTradeDay": lastLastTradeDay.Format(global.ShortTimeLayout),
	}).Info("Last Trade Days")
}
