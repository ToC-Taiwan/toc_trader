// Package sysparm package sysparm
package sysparm

import (
	"gorm.io/gorm"
)

// Parameters Parameters
type Parameters struct {
	gorm.Model `json:"-" swaggerignore:"true"`
	Key        string `gorm:"column:key;uniqueIndex" json:"key"`
	Value      string `gorm:"column:value" json:"value"`
}

// TargetCondArr TargetCondArr
type TargetCondArr struct {
	LimitPriceLow   float64 `json:"limit_price_low"`
	LimitPriceHigh  float64 `json:"limit_price_high"`
	LimitVolumeLow  int64   `json:"limit_volume_low"`
	LimitVolumeHigh int64   `json:"limit_volume_high"`
}

// Tabler Tabler
type Tabler interface {
	TableName() string
}

// TableName TableName
func (Parameters) TableName() string {
	return "settings"
}
