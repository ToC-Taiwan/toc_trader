// Package simulateprocess package simulateprocess
package simulateprocess

import (
	"time"

	"gitlab.tocraw.com/root/toc_trader/pkg/global"
)

func getLastTradeOutTimeByEntireTickTimeStamp(timeStamp int64) int64 {
	tmp := time.Unix(0, timeStamp)
	endTime := time.Date(tmp.Year(), tmp.Month(), tmp.Day(), global.TradeOutEndHour, global.TradeOutEndMinute, 0, 0, time.UTC)
	return endTime.UnixNano()
}

func getLastTradeInTimeByEntireTickTimeStamp(timeStamp int64) int64 {
	tmp := time.Unix(0, timeStamp)
	endTime := time.Date(tmp.Year(), tmp.Month(), tmp.Day(), global.TradeInEndHour, global.TradeInEndMinute, 0, 0, time.UTC)
	return endTime.UnixNano()
}

type simulateType int

const (
	simTypeForward simulateType = iota + 1
	simTypeReverse
)
