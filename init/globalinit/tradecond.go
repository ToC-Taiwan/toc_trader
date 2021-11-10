// Package globalinit is init all global var
package globalinit

import (
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
)

func init() {
	// Forward Condition
	global.ForwardCond = simulationcond.AnalyzeCondition{
		HistoryCloseCount:     2100,
		TrimHistoryCloseCount: true,
		OutInRatio:            85,
		ReverseOutInRatio:     0,
		CloseChangeRatioHigh:  3,
		CloseChangeRatioLow:   -1,
		CloseDiff:             0,
		OpenChangeRatio:       -1,
		ReverseRsiHigh:        0.8,
		ReverseRsiLow:         0.1,
		RsiHigh:               0.8,
		RsiLow:                0.1,
		TicksPeriodCount:      2,
		TicksPeriodLimit:      4 * 1.3,
		TicksPeriodThreshold:  4,
		VolumePerSecond:       30,
		MaxHoldTime:           1,
	}
	// Reverse Condition
	global.ReverseCond = simulationcond.AnalyzeCondition{
		HistoryCloseCount:     1100,
		TrimHistoryCloseCount: true,
		OutInRatio:            100,
		ReverseOutInRatio:     6,
		CloseChangeRatioHigh:  3,
		CloseChangeRatioLow:   0,
		CloseDiff:             0,
		OpenChangeRatio:       3,
		ReverseRsiHigh:        0.7,
		ReverseRsiLow:         0.1,
		RsiHigh:               0.7,
		RsiLow:                0.1,
		TicksPeriodCount:      3,
		TicksPeriodLimit:      8 * 1.3,
		TicksPeriodThreshold:  8,
		VolumePerSecond:       10,
		MaxHoldTime:           1,
	}
}
