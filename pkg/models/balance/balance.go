// Package balance package balance
package balance

import (
	"time"

	"gorm.io/gorm"
)

// Balance Balance
type Balance struct {
	gorm.Model
	TradeDay        time.Time `gorm:"column:trade_day;"`
	TradeCount      int64     `gorm:"column:trade_count;"`
	Forward         int64     `gorm:"column:forward"`
	Reverse         int64     `gorm:"column:reverse"`
	OriginalBalance int64     `gorm:"column:original_balance"`
	Discount        int64     `gorm:"column:discount"`
	Total           int64     `gorm:"column:total"`
}

// Tabler Tabler
type Tabler interface {
	TableName() string
}

// TableName TableName
func (Balance) TableName() string {
	return "trade_balance"
}
