// Package simulateprocess package simulateprocess
package simulateprocess

import (
	"time"
)

func getLastTradeTimeByEntireTickTimeStamp(timeStamp int64) int64 {
	tmp := time.Unix(0, timeStamp)
	endTime := time.Date(tmp.Year(), tmp.Month(), tmp.Day(), 13, 0, 0, 0, time.UTC)
	return endTime.UnixNano()
}
