// Package holiday package holiday
package holiday

import "gorm.io/gorm"

// Holiday Holiday
type Holiday struct {
	gorm.Model `json:"-" swaggerignore:"true"`
	TimeStamp  int64 `gorm:"column:timestamp;uniqueIndex"`
}

// Tabler Tabler
type Tabler interface {
	TableName() string
}

// TableName TableName
func (Holiday) TableName() string {
	return "basic_holiday"
}
