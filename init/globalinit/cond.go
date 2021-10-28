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
		HistoryCloseCount:     1500,
		OutInRatio:            90,
		ReverseOutInRatio:     10,
		CloseDiff:             0,
		CloseChangeRatioLow:   -3,
		CloseChangeRatioHigh:  7,
		OpenChangeRatio:       7,
		RsiHigh:               50.9,
		RsiLow:                50.3,
		ReverseRsiHigh:        50.9,
		ReverseRsiLow:         50.3,
		TicksPeriodThreshold:  8,
		TicksPeriodLimit:      8 * 1.3,
		TicksPeriodCount:      1,
		VolumePerSecond:       12,
	}

	global.ReverseCond = simulationcond.AnalyzeCondition{
		TrimHistoryCloseCount: false,
		HistoryCloseCount:     1500,
		OutInRatio:            90,
		ReverseOutInRatio:     10,
		CloseDiff:             0,
		CloseChangeRatioLow:   -3,
		CloseChangeRatioHigh:  7,
		OpenChangeRatio:       7,
		RsiHigh:               49.2,
		RsiLow:                48.8,
		ReverseRsiHigh:        49.2,
		ReverseRsiLow:         48.8,
		TicksPeriodThreshold:  8,
		TicksPeriodLimit:      8 * 1.3,
		TicksPeriodCount:      1,
		VolumePerSecond:       12,
	}
}
