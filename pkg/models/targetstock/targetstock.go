// Package targetstock package targetstock
package targetstock

import (
	"time"

	"gitlab.tocraw.com/root/toc_trader/pkg/models/stock"
	"gorm.io/gorm"
)

// Target Target
type Target struct {
	gorm.Model
	LastTradeDay time.Time   `gorm:"column:last_trade_day"`
	Stock        stock.Stock `gorm:"foreignKey:StockID"`
	StockID      int64
}

// TableName TableName
func (Target) TableName() string {
	return "basic_stock_target"
}
