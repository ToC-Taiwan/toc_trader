// Package simulate package simulate
package simulate

import (
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
	"gorm.io/gorm"
)

// Result Result
type Result struct {
	gorm.Model
	Balance        int64 `gorm:"column:balance;"`
	ForwardBalance int64 `gorm:"column:forward_balance;"`
	ReverseBalance int64 `gorm:"column:reverse_balance;"`
	TotalLoss      int64 `gorm:"column:total_loss;"`
	TradeCount     int64 `gorm:"column:trade_count;"`
	PositiveDays   int64 `gorm:"column:positive_days;"`
	TotalDays      int64 `gorm:"column:total_days;"`
	IsBestForward  bool  `gorm:"column:is_best_forward;"`
	IsBestReverse  bool  `gorm:"column:is_best_reverse;"`
	CondID         int64
	Cond           simulationcond.AnalyzeCondition `gorm:"foreignKey:CondID"`
}

// Tabler Tabler
type Tabler interface {
	TableName() string
}

// TableName TableName
func (Result) TableName() string {
	return "simulate_result"
}
