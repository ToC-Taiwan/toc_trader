// Package globalinit is init all global var
package globalinit

import (
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
)

func init() {
	global.CentralCond = simulationcond.AnalyzeCondition{
		TrimHistoryCloseCount: true,
		HistoryCloseCount:     1500,
		OutInRatio:            60,
		ReverseOutInRatio:     10,
		CloseDiff:             0,
		CloseChangeRatioLow:   -3,
		CloseChangeRatioHigh:  7,
		OpenChangeRatio:       7,
		RsiHigh:               50,
		RsiLow:                50,
		ReverseRsiHigh:        50,
		ReverseRsiLow:         50,
		TicksPeriodThreshold:  10,
		TicksPeriodLimit:      10 * 1.3,
		TicksPeriodCount:      1,
		VolumePerSecond:       5,
	}

	global.ForwardCond = simulationcond.AnalyzeCondition{
		TrimHistoryCloseCount: false,
		HistoryCloseCount:     2500,
		OutInRatio:            90,
		ReverseOutInRatio:     3,
		CloseDiff:             0,
		CloseChangeRatioLow:   -3,
		CloseChangeRatioHigh:  7,
		OpenChangeRatio:       7,
		RsiHigh:               50.8,
		RsiLow:                50.6,
		ReverseRsiHigh:        50.8,
		ReverseRsiLow:         50.6,
		TicksPeriodThreshold:  6,
		TicksPeriodLimit:      6 * 1.3,
		TicksPeriodCount:      2,
		VolumePerSecond:       12,
	}

	global.ReverseCond = simulationcond.AnalyzeCondition{
		TrimHistoryCloseCount: false,
		HistoryCloseCount:     2000,
		OutInRatio:            90,
		ReverseOutInRatio:     3,
		CloseDiff:             0,
		CloseChangeRatioLow:   -3,
		CloseChangeRatioHigh:  7,
		OpenChangeRatio:       7,
		RsiHigh:               49.5,
		RsiLow:                48.7,
		ReverseRsiHigh:        49.5,
		ReverseRsiLow:         48.7,
		TicksPeriodThreshold:  10,
		TicksPeriodLimit:      10 * 1.3,
		TicksPeriodCount:      1,
		VolumePerSecond:       8,
	}
}
