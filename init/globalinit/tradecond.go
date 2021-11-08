// Package globalinit is init all global var
package globalinit

import (
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
)

func init() {
	// Forward Condition
	global.ForwardCond = simulationcond.AnalyzeCondition{
		HistoryCloseCount:     2000,
		TrimHistoryCloseCount: true,
		OutInRatio:            90,
		ReverseOutInRatio:     0,
		CloseChangeRatioHigh:  3,
		CloseChangeRatioLow:   -1,
		CloseDiff:             0,
		OpenChangeRatio:       -1,
		ReverseRsiHigh:        0.9,
		ReverseRsiLow:         0.1,
		RsiHigh:               0.9,
		RsiLow:                0.1,
		TicksPeriodCount:      3,
		TicksPeriodLimit:      4 * 1.3,
		TicksPeriodThreshold:  4,
		VolumePerSecond:       30,
	}
	// Reverse Condition
	global.ReverseCond = simulationcond.AnalyzeCondition{
		HistoryCloseCount:     2300,
		TrimHistoryCloseCount: true,
		OutInRatio:            100,
		ReverseOutInRatio:     9,
		CloseChangeRatioHigh:  3,
		CloseChangeRatioLow:   0,
		CloseDiff:             0,
		OpenChangeRatio:       3,
		ReverseRsiHigh:        0.9,
		ReverseRsiLow:         0.2,
		RsiHigh:               0.9,
		RsiLow:                0.2,
		TicksPeriodCount:      1,
		TicksPeriodLimit:      8 * 1.3,
		TicksPeriodThreshold:  8,
		VolumePerSecond:       30,
	}
}
