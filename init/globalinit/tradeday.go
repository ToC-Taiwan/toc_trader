// Package globalinit is init all global var
package globalinit

import (
	"time"

	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/importbasic"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
)

func init() {
	var err error
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
