// Package globalinit is init all global var
package globalinit

import (
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
)

func init() {
	global.BaseCond = simulationcond.AnalyzeCondition{
		TrimHistoryCloseCount: true,
		HistoryCloseCount:     600,
		OutInRatio:            85,
		ReverseOutInRatio:     3,
		CloseDiff:             0,
		CloseChangeRatioLow:   0,
		CloseChangeRatioHigh:  3,
		OpenChangeRatio:       3,
		RsiHigh:               0.8,
		RsiLow:                0.2,
		ReverseRsiHigh:        0.8,
		ReverseRsiLow:         0.2,
		TicksPeriodThreshold:  4,
		TicksPeriodLimit:      4 * 1.3,
		TicksPeriodCount:      2,
		VolumePerSecond:       10,
	}
}
