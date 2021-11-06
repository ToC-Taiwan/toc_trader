// Package globalinit is init all global var
package globalinit

import (
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
)

func init() {
	// Forward Condition
	global.ForwardCond = simulationcond.AnalyzeCondition{
		TrimHistoryCloseCount: true,
		HistoryCloseCount:     800,
		OutInRatio:            95,
		ReverseOutInRatio:     0,
		CloseDiff:             0,
		CloseChangeRatioLow:   -1,
		CloseChangeRatioHigh:  3,
		OpenChangeRatio:       -1,
		RsiHigh:               0.9,
		RsiLow:                0.1,
		ReverseRsiHigh:        0.9,
		ReverseRsiLow:         0.1,
		TicksPeriodThreshold:  8,
		TicksPeriodLimit:      8 * 1.3,
		TicksPeriodCount:      3,
		VolumePerSecond:       20,
	}
	// Reverse Condition
	global.ReverseCond = simulationcond.AnalyzeCondition{
		TrimHistoryCloseCount: true,
		HistoryCloseCount:     2100,
		OutInRatio:            100,
		ReverseOutInRatio:     9,
		CloseDiff:             0,
		CloseChangeRatioLow:   0,
		CloseChangeRatioHigh:  3,
		OpenChangeRatio:       3,
		RsiHigh:               0.7,
		RsiLow:                0.1,
		ReverseRsiHigh:        0.7,
		ReverseRsiLow:         0.1,
		TicksPeriodThreshold:  4,
		TicksPeriodLimit:      4 * 1.3,
		TicksPeriodCount:      2,
		VolumePerSecond:       30,
	}
}
