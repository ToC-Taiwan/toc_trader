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
	global.MeanTimeTradeStockNum = 8

	global.TickAnalyzeCondition = global.AnalyzeCondition{
		HistoryCloseCount:    400,
		OutSum:               300,
		OutInRatio:           55,
		CloseDiff:            0,
		CloseChangeRatioLow:  -2,
		CloseChangeRatioHigh: 7,
		OpenChangeRatio:      7,
		RsiHigh:              50,
		RsiLow:               50,
	}

	if err := importbasic.ImportHoliday(); err != nil {
		panic(err)
	}

	var today time.Time
	var err error
	if time.Now().Hour() >= 15 {
		today = time.Now().AddDate(0, 0, 1)
	} else {
		today = time.Now()
	}
	global.TradeDay, err = importbasic.GetTradeDayTime(today)
	if err != nil {
		panic(err)
	}

	global.TradeDayEndTime = time.Date(global.TradeDay.Year(), global.TradeDay.Month(), global.TradeDay.Day(), 13, 0, 0, 0, time.Local)

	global.LastTradeDay, err = importbasic.GetLastTradeDayTime(global.TradeDay)
	if err != nil {
		panic(err)
	}

	global.LastLastTradeDay, err = importbasic.GetLastTradeDayTime(global.LastTradeDay)
	if err != nil {
		panic(err)
	}

	global.LastTradeDayArr = append(global.LastTradeDayArr, global.LastLastTradeDay, global.LastTradeDay)

	logger.Logger.WithFields(map[string]interface{}{
		"TradeDay":         global.TradeDay.Format(global.ShortTimeLayout),
		"LastTradeDay":     global.LastTradeDay.Format(global.ShortTimeLayout),
		"LastLastTradeDay": global.LastLastTradeDay.Format(global.ShortTimeLayout),
	}).Info("Last Trade Days")
	logger.Logger.WithFields(map[string]interface{}{
		"TradeDayEndTime": global.TradeDayEndTime.Format(global.LongTimeLayout),
	}).Info("Trade Time")
}
