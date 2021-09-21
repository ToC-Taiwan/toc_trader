// Package traderole package traderole
package traderole

import "gorm.io/gorm"

// Role Role
type Role struct {
	gorm.Model
	Title    string `gorm:"column:title;uniqueIndex"`
	IsActive bool   `gorm:"column:is_active;default:true"`
}

// TableName TableName
func (Role) TableName() string {
	return "basic_role"
}
