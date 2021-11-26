// Package stock package stock
package stock

import (
	"sync"

	"gorm.io/gorm"
)

// Stock Stock
type Stock struct {
	gorm.Model `json:"-" swaggerignore:"true"`
	StockNum   string  `gorm:"column:stock_num;uniqueIndex;index:idx_stock"`
	Category   string  `gorm:"column:category"`
	StockName  string  `gorm:"column:stock_name"`
	StockType  string  `gorm:"column:stock_type"`
	DayTrade   bool    `gorm:"column:day_trade;index:idx_stock"`
	LastClose  float64 `gorm:"column:last_close"`
	// LastVolume int64   `gorm:"column:last_volume"`
}

// Tabler Tabler
type Tabler interface {
	TableName() string
}

// TableName TableName
func (Stock) TableName() string {
	return "basic_stock"
}

// MutexStruct MutexStruct
type MutexStruct struct {
	dataMap map[string]Stock
	mutex   sync.RWMutex
}

// Set Set
func (c *MutexStruct) Set(stock Stock) {
	if c.dataMap == nil {
		c.dataMap = make(map[string]Stock)
	}
	c.mutex.Lock()
	c.dataMap[stock.StockNum] = stock
	c.mutex.Unlock()
}

// Get Get
func (c *MutexStruct) Get(stockNum string) Stock {
	var tmp Stock
	c.mutex.RLock()
	tmp = c.dataMap[stockNum]
	c.mutex.RUnlock()
	return tmp
}

// GetCategory GetCategory
func (c *MutexStruct) GetCategory(stockNum string) string {
	var tmp Stock
	c.mutex.RLock()
	tmp = c.dataMap[stockNum]
	c.mutex.RUnlock()
	return tmp.Category
}

// CheckIsDayTrade CheckIsDayTrade
func (c *MutexStruct) CheckIsDayTrade(stockNum string) bool {
	var tmp bool
	c.mutex.RLock()
	tmp = c.dataMap[stockNum].DayTrade
	c.mutex.RUnlock()
	return tmp
}
