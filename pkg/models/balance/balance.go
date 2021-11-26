// Package balance package balance
package balance

import (
	"time"

	"gorm.io/gorm"
)

// Balance Balance
type Balance struct {
	gorm.Model      `json:"-" swaggerignore:"true"`
	TradeDay        time.Time `gorm:"column:trade_day;" json:"trade_day"`
	TradeCount      int64     `gorm:"column:trade_count;" json:"trade_count"`
	Forward         int64     `gorm:"column:forward" json:"forward"`
	Reverse         int64     `gorm:"column:reverse" json:"reverse"`
	OriginalBalance int64     `gorm:"column:original_balance" json:"original_balance"`
	Discount        int64     `gorm:"column:discount" json:"discount"`
	Total           int64     `gorm:"column:total" json:"total"`
}

// Tabler Tabler
type Tabler interface {
	TableName() string
}

// TableName TableName
func (Balance) TableName() string {
	return "trade_balance"
}
