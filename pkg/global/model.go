// Package global package global
package global

// AnalyzeCondition AnalyzeCondition
type AnalyzeCondition struct {
	HistoryCloseCount    int64
	OutSum               int64
	OutInRatio           float64
	CloseDiff            float64
	CloseChangeRatioLow  float64
	CloseChangeRatioHigh float64
	OpenChangeRatio      float64
	RsiHigh              int64
	RsiLow               int64
	TicksPeriodThreshold float64
	TicksPeriodLimit     float64
	TicksPeriodCount     int
}

// SystemSwitch SystemSwitch
type SystemSwitch struct {
	Buy                          bool
	Sell                         bool
	SellFirst                    bool
	BuyLater                     bool
	UseBidAsk                    bool
	MeanTimeTradeStockNum        int
	MeanTimeReverseTradeStockNum int
}
