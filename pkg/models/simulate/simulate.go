// Package simulate package simulate
package simulate

import (
	"time"

	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
	"gorm.io/gorm"
)

// Result Result
type Result struct {
	gorm.Model     `json:"-" swaggerignore:"true"`
	Balance        int64                           `gorm:"column:balance;" json:"balance"`
	ForwardBalance int64                           `gorm:"column:forward_balance;" json:"forward_balance"`
	ReverseBalance int64                           `gorm:"column:reverse_balance;" json:"reverse_balance"`
	TotalLoss      int64                           `gorm:"column:total_loss;" json:"total_loss"`
	TradeCount     int64                           `gorm:"column:trade_count;" json:"trade_count"`
	PositiveDays   int64                           `gorm:"column:positive_days;" json:"positive_days"`
	NegativeDays   int64                           `gorm:"column:negative_days;" json:"negative_days"`
	TotalDays      int64                           `gorm:"column:total_days;" json:"total_days"`
	IsBestForward  bool                            `gorm:"column:is_best_forward;" json:"is_best_forward"`
	IsBestReverse  bool                            `gorm:"column:is_best_reverse;" json:"is_best_reverse"`
	TradeDay       time.Time                       `gorm:"column:trade_day;" json:"trade_day"`
	CondID         int64                           `gorm:"column:cond_id;" json:"cond_id"`
	Cond           simulationcond.AnalyzeCondition `gorm:"foreignKey:CondID" json:"cond"`
}

// Tabler Tabler
type Tabler interface {
	TableName() string
}

// TableName TableName
func (Result) TableName() string {
	return "simulate_result"
}
