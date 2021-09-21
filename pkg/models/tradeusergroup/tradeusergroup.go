// Package tradeusergroup package tradeusergroup
package tradeusergroup

import "gorm.io/gorm"

// UserGroup UserGroup
type UserGroup struct {
	gorm.Model
	Title string `gorm:"column:title"`
}

// Tabler Tabler
type Tabler interface {
	TableName() string
}

// TableName TableName
func (UserGroup) TableName() string {
	return "basic_user_group"
}
