// Package tradeevent package tradeevent
package tradeevent

import "gorm.io/gorm"

// EventResponse EventResponse
type EventResponse struct {
	gorm.Model `json:"-" swaggerignore:"true"`
	Event      string `gorm:"column:event" json:"event"`
	EventCode  int64  `gorm:"column:event_code" json:"event_code"`
	Info       string `gorm:"column:info" json:"info"`
	Response   int64  `gorm:"column:response" json:"response"`
}

// Tabler Tabler
type Tabler interface {
	TableName() string
}

// TableName TableName
func (EventResponse) TableName() string {
	return "trade_event"
}
