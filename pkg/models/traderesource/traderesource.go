// Package traderesource package traderesource
package traderesource

import "gorm.io/gorm"

// Resource Resource
type Resource struct {
	gorm.Model
	Title string `gorm:"column:title;uniqueIndex"`
}

// TableName TableName
func (Resource) TableName() string {
	return "basic_resource"
}
