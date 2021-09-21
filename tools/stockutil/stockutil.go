// Package stockutil is Utils for stocks
package stockutil

import (
	"gitlab.tocraw.com/root/toc_trader/tools/common"
)

// GetNewClose GetNewClose
func GetNewClose(close float64, unit int64) float64 {
	if close == 0 {
		return 0
	}
	for {
		if unit == 0 {
			return common.Round(close, 2)
		}
		diff := GetDiff(close)
		if unit > 0 {
			close += diff
			unit--
		} else {
			close -= diff
			unit++
		}
	}
}

// GetDiff GetDiff
func GetDiff(close float64) float64 {
	switch {
	case close > 0 && close < 10:
		return 0.01
	case close >= 10 && close < 50:
		return 0.05
	case close >= 50 && close < 100:
		return 0.1
	case close >= 100 && close < 500:
		return 0.5
	case close >= 500 && close < 1000:
		return 1
	case close >= 1000:
		return 5
	}
	return 0
}
