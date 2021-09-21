// Package analyzeentiretick analyzeentiretick
package analyzeentiretick

import (
	"time"

	"gitlab.tocraw.com/root/toc_trader/pkg/models/analyzestreamtick"
	"gorm.io/gorm"
)

// AnalyzeEntireTick AnalyzeEntireTick
type AnalyzeEntireTick struct {
	gorm.Model
	TimeStamp        int64   `gorm:"column:timestamp;index:idx_analyzeentiretick"`
	StockNum         string  `gorm:"column:stock_num;index:idx_analyzeentiretick"`
	Close            float64 `gorm:"column:close"`
	Rsi              float64 `gorm:"column:rsi"`
	CloseChangeRatio float64 `gorm:"column:close_change_ratio"`
	OpenChangeRatio  float64 `gorm:"column:open_change_ratio"`
	OutSum           int64   `gorm:"column:out_sum"`
	InSum            int64   `gorm:"column:in_sum"`
	OutInRatio       float64 `gorm:"column:out_in_ratio"`
	TotalTime        float64 `gorm:"column:total_time"`
	CloseDiff        float64 `gorm:"column:close_diff"`
	Open             float64 `gorm:"column:open"`
	AvgPrice         float64 `gorm:"column:avg_price"`
	High             float64 `gorm:"column:high"`
	Low              float64 `gorm:"column:low"`
}

// Tabler Tabler
type Tabler interface {
	TableName() string
}

// TableName TableName
func (AnalyzeEntireTick) TableName() string {
	return "analyze_entire_tick"
}

// ToAnalyzeStreamTick ToAnalyzeStreamTick
func (c *AnalyzeEntireTick) ToAnalyzeStreamTick() *analyzestreamtick.AnalyzeStreamTick {
	tickTime := time.Unix(0, c.TimeStamp).Add(-8 * time.Hour).UnixNano()
	tmp := analyzestreamtick.AnalyzeStreamTick{
		TimeStamp:        tickTime,
		StockNum:         c.StockNum,
		Close:            c.Close,
		Rsi:              c.Rsi,
		OpenChangeRatio:  c.OpenChangeRatio,
		CloseChangeRatio: c.CloseChangeRatio,
		OutSum:           c.OutSum,
		InSum:            c.InSum,
		OutInRatio:       c.OutInRatio,
		TotalTime:        c.TotalTime,
		CloseDiff:        c.CloseDiff,
		Open:             c.Open,
		AvgPrice:         c.AvgPrice,
		High:             c.High,
		Low:              c.Low,
	}
	return &tmp
}
