// Package simulate package simulate
package simulate

import (
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
	"gorm.io/gorm"
)

// Result Result
type Result struct {
	gorm.Model
	Balance int64                           `gorm:"column:balance;"`
	Cond    simulationcond.AnalyzeCondition `gorm:"foreignKey:CondID"`
	CondID  int64
}

// Tabler Tabler
type Tabler interface {
	TableName() string
}

// TableName TableName
func (Result) TableName() string {
	return "simulate_result"
}
