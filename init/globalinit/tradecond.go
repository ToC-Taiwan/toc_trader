// Package globalinit is init all global var
package globalinit

import (
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
)

func init() {
	// Forward Condition
	global.ForwardCond = simulationcond.AnalyzeCondition{
		HistoryCloseCount:     2300,
		TrimHistoryCloseCount: true,
		ForwardOutInRatio:     85,
		CloseChangeRatioLow:   -1,
		CloseChangeRatioHigh:  3,
		OpenChangeRatio:       -1,
		RsiHigh:               0.8,
		TicksPeriodThreshold:  12,
		TicksPeriodLimit:      12 * 1.3,
		TicksPeriodCount:      2,
		VolumePerSecond:       30,
		MaxHoldTime:           3,
	}
	// Reverse Condition
	global.ReverseCond = simulationcond.AnalyzeCondition{
		HistoryCloseCount:     1800,
		TrimHistoryCloseCount: true,
		ReverseOutInRatio:     3,
		CloseChangeRatioLow:   0,
		CloseChangeRatioHigh:  3,
		OpenChangeRatio:       3,
		RsiLow:                0.2,
		TicksPeriodThreshold:  4,
		TicksPeriodLimit:      4 * 1.3,
		TicksPeriodCount:      3,
		VolumePerSecond:       10,
		MaxHoldTime:           1,
	}
}
