// Package targetstock package targetstock
package targetstock

import (
	"time"

	"gitlab.tocraw.com/root/toc_trader/pkg/models/stock"
	"gorm.io/gorm"
)

// Target Target
type Target struct {
	gorm.Model `json:"-" swaggerignore:"true"`
	TradeDay   time.Time   `gorm:"column:trade_day"`
	Rank       int64       `gorm:"column:rank"`
	StockID    int64       `gorm:"column:stock_id"`
	Stock      stock.Stock `gorm:"foreignKey:StockID"`
}

// TableName TableName
func (Target) TableName() string {
	return "basic_stock_target"
}
