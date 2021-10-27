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
		OutInRatio:            80,
		ReverseOutInRatio:     3,
		CloseDiff:             0,
		CloseChangeRatioLow:   -3,
		CloseChangeRatioHigh:  7,
		OpenChangeRatio:       7,
		RsiHigh:               51.1,
		RsiLow:                50.7,
		ReverseRsiHigh:        51.1,
		ReverseRsiLow:         50.7,
		TicksPeriodThreshold:  8,
		TicksPeriodLimit:      8 * 1.3,
		TicksPeriodCount:      1,
		VolumePerSecond:       6,
	}

	global.ReverseCond = simulationcond.AnalyzeCondition{
		TrimHistoryCloseCount: false,
		HistoryCloseCount:     2000,
		OutInRatio:            80,
		ReverseOutInRatio:     9,
		CloseDiff:             0,
		CloseChangeRatioLow:   -3,
		CloseChangeRatioHigh:  7,
		OpenChangeRatio:       7,
		RsiHigh:               49.8,
		RsiLow:                49,
		ReverseRsiHigh:        49.8,
		ReverseRsiLow:         49,
		TicksPeriodThreshold:  8,
		TicksPeriodLimit:      8 * 1.3,
		TicksPeriodCount:      1,
		VolumePerSecond:       12,
	}
}
