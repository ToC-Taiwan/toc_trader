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
}
