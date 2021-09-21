// Package tradeuser package tradeuser
package tradeuser

import "gorm.io/gorm"

// User User
type User struct {
	gorm.Model
	Username string `gorm:"column:username;uniqueIndex"`
	Password string `gorm:"column:password"`
	IsActive bool   `gorm:"column:is_active;default:true"`
	IsSuper  bool   `gorm:"column:is_super;default:false"`
	Group    int64  `gorm:"column:group"`
}

// Tabler Tabler
type Tabler interface {
	TableName() string
}

// TableName TableName
func (User) TableName() string {
	return "basic_user"
}
