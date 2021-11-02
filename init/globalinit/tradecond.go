// Package globalinit is init all global var
package globalinit

import (
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
)

func init() {
	global.ForwardCond = simulationcond.AnalyzeCondition{
		TrimHistoryCloseCount: true,
		HistoryCloseCount:     1000,
		OutInRatio:            85,
		ReverseOutInRatio:     3,
		CloseDiff:             0,
		CloseChangeRatioLow:   0,
		CloseChangeRatioHigh:  3,
		OpenChangeRatio:       3,
		RsiHigh:               0.7,
		RsiLow:                0.2,
		ReverseRsiHigh:        0.7,
		ReverseRsiLow:         0.2,
		TicksPeriodThreshold:  4,
		TicksPeriodLimit:      4 * 1.3,
		TicksPeriodCount:      1,
		VolumePerSecond:       20,
	}

	global.ReverseCond = simulationcond.AnalyzeCondition{
		TrimHistoryCloseCount: true,
		HistoryCloseCount:     700,
		OutInRatio:            85,
		ReverseOutInRatio:     3,
		CloseDiff:             0,
		CloseChangeRatioLow:   0,
		CloseChangeRatioHigh:  3,
		OpenChangeRatio:       3,
		RsiHigh:               0.9,
		RsiLow:                0.1,
		ReverseRsiHigh:        0.9,
		ReverseRsiLow:         0.1,
		TicksPeriodThreshold:  4,
		TicksPeriodLimit:      4 * 1.3,
		TicksPeriodCount:      2,
		VolumePerSecond:       30,
	}
}
