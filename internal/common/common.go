// Package common is common tools
package common

import (
	"math"
	"strconv"
	"time"

	"gitlab.tocraw.com/root/toc_trader/pkg/global"
)

// StrToInt64 StrToInt64
func StrToInt64(input string) (ans int64, err error) {
	ans, err = strconv.ParseInt(input, 10, 64)
	if err != nil {
		return ans, err
	}
	return ans, err
}

// StrToFloat64 StrToFloat64
func StrToFloat64(input string) (ans float64, err error) {
	ans, err = strconv.ParseFloat(input, 64)
	if err != nil {
		return ans, err
	}
	return ans, err
}

// Round Round
func Round(val float64, precision int) float64 {
	p := math.Pow10(precision)
	return math.Floor(val*p+0.5) / p
}

// MicroDateTimeToTimeStamp MicroDateTimeToTimeStamp
func MicroDateTimeToTimeStamp(dateTime string) (timeStamp int64, err error) {
	dataTimeStamp, err := time.ParseInLocation(global.LongTimeLayout, dateTime[:19], time.Local)
	if err != nil {
		return timeStamp, err
	}
	microSec, err := StrToInt64(dateTime[20:])
	if err != nil {
		return timeStamp, err
	}
	return dataTimeStamp.UnixNano() + microSec*1000, err
}
