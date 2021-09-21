// Package analyzestreamtick analyzestreamtick
package analyzestreamtick

import (
	"sync"

	"gorm.io/gorm"
)

// AnalyzeStreamTick AnalyzeStreamTick
type AnalyzeStreamTick struct {
	gorm.Model
	TimeStamp        int64   `gorm:"column:timestamp;index:idx_analyzestreamtick"`
	StockNum         string  `gorm:"column:stock_num;index:idx_analyzestreamtick"`
	Close            float64 `gorm:"column:close"`
	Rsi              float64 `gorm:"column:rsi"`
	OpenChangeRatio  float64 `gorm:"column:open_change_ratio"`
	CloseChangeRatio float64 `gorm:"column:close_change_ratio"`
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
func (AnalyzeStreamTick) TableName() string {
	return "analyze_stream_tick"
}

// AnalyzeStreamArrMutexStruct AnalyzeStreamArrMutexStruct
type AnalyzeStreamArrMutexStruct struct {
	Ticks []*AnalyzeStreamTick
	Mutex sync.RWMutex
}

// Append Append
func (c *AnalyzeStreamArrMutexStruct) Append(value *AnalyzeStreamTick) {
	c.Mutex.Lock()
	c.Ticks = append(c.Ticks, value)
	c.Mutex.Unlock()
}

// Clear Clear
func (c *AnalyzeStreamArrMutexStruct) Clear() {
	c.Mutex.Lock()
	c.Ticks = []*AnalyzeStreamTick{}
	c.Mutex.Unlock()
}

// GetTotalCount GetTotalCount
func (c *AnalyzeStreamArrMutexStruct) GetTotalCount() int {
	var tmp int
	c.Mutex.RLock()
	tmp = len(c.Ticks)
	c.Mutex.RUnlock()
	return tmp
}
