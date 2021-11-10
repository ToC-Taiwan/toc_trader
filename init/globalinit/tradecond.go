// Package globalinit is init all global var
package globalinit

import (
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
)

func init() {
	// Forward Condition
	global.ForwardCond = simulationcond.AnalyzeCondition{
		HistoryCloseCount:     1800,
		TrimHistoryCloseCount: true,
		OutInRatio:            95,
		ReverseOutInRatio:     0,
		CloseChangeRatioHigh:  3,
		CloseChangeRatioLow:   -1,
		CloseDiff:             0,
		OpenChangeRatio:       -1,
		ReverseRsiHigh:        0.9,
		ReverseRsiLow:         0.1,
		RsiHigh:               0.9,
		RsiLow:                0.1,
		TicksPeriodCount:      2,
		TicksPeriodLimit:      8 * 1.3,
		TicksPeriodThreshold:  8,
		VolumePerSecond:       10,
	}
	// Reverse Condition
	global.ReverseCond = simulationcond.AnalyzeCondition{
		HistoryCloseCount:     1900,
		TrimHistoryCloseCount: true,
		OutInRatio:            100,
		ReverseOutInRatio:     3,
		CloseChangeRatioHigh:  3,
		CloseChangeRatioLow:   0,
		CloseDiff:             0,
		OpenChangeRatio:       3,
		ReverseRsiHigh:        0.7,
		ReverseRsiLow:         0.1,
		RsiHigh:               0.7,
		RsiLow:                0.1,
		TicksPeriodCount:      2,
		TicksPeriodLimit:      12 * 1.3,
		TicksPeriodThreshold:  12,
		VolumePerSecond:       10,
	}
}
