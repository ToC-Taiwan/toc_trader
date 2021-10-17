// Package simulationcond package simulationcond
package simulationcond

import (
	"gorm.io/gorm"
)

// AnalyzeCondition AnalyzeCondition
type AnalyzeCondition struct {
	gorm.Model
	HistoryCloseCount    int64   `gorm:"column:history_close_count"`
	OutInRatio           float64 `gorm:"column:out_in_ratio"`
	ReverseOutInRatio    float64 `gorm:"column:reverse_out_in_ratio"`
	CloseDiff            float64 `gorm:"column:close_diff"`
	CloseChangeRatioLow  float64 `gorm:"column:close_change_ratio_low"`
	CloseChangeRatioHigh float64 `gorm:"column:close_change_ratio_high"`
	OpenChangeRatio      float64 `gorm:"column:open_change_ratio"`
	RsiHigh              float64 `gorm:"column:rsi_high"`
	RsiLow               float64 `gorm:"column:rsi_low"`
	ReverseRsiHigh       float64 `gorm:"column:reverse_rsi_high"`
	ReverseRsiLow        float64 `gorm:"column:reverse_rsi_low"`
	TicksPeriodThreshold float64 `gorm:"column:ticks_period_threshold"`
	TicksPeriodLimit     float64 `gorm:"column:ticks_period_limit"`
	TicksPeriodCount     int     `gorm:"column:ticks_period_count"`
	VolumePerSecond      int64   `gorm:"column:volume_per_second"`
}

// Tabler Tabler
type Tabler interface {
	TableName() string
}

// TableName TableName
func (AnalyzeCondition) TableName() string {
	return "simulate_cond"
}
