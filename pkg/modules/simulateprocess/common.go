// Package simulateprocess package simulateprocess
package simulateprocess

import (
	"time"

	"gitlab.tocraw.com/root/toc_trader/pkg/global"
)

func getLastTradeTimeByEntireTickTimeStamp(timeStamp int64) int64 {
	tmp := time.Unix(0, timeStamp)
	endTime := time.Date(tmp.Year(), tmp.Month(), tmp.Day(), global.TradeEndHour, global.TradeEndMinute, 0, 0, time.UTC)
	return endTime.UnixNano()
}
