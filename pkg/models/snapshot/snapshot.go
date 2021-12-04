// Package snapshot package snapshot
package snapshot

import "gorm.io/gorm"

// SnapShot SnapShot
type SnapShot struct {
	gorm.Model      `json:"-" swaggerignore:"true"`
	TS              int64   `gorm:"column:ts"`
	Code            string  `gorm:"column:code"`
	Exchange        string  `gorm:"column:exchange"`
	Open            float64 `gorm:"column:open"`
	High            float64 `gorm:"column:high"`
	Low             float64 `gorm:"column:low"`
	Close           float64 `gorm:"column:close"`
	TickType        string  `gorm:"column:tick_type"`
	ChangePrice     float64 `gorm:"column:change_price"`
	ChangeRate      float64 `gorm:"column:change_rate"`
	ChangeType      string  `gorm:"column:change_type"`
	AveragePrice    float64 `gorm:"column:average_price"`
	Volume          int64   `gorm:"column:volume"`
	TotalVolume     int64   `gorm:"column:total_volume"`
	Amount          int64   `gorm:"column:amount"`
	TotalAmount     int64   `gorm:"column:total_amount"`
	YesterdayVolume float64 `gorm:"column:yesterday_volume"`
	BuyPrice        float64 `gorm:"column:buy_price"`
	BuyVolume       float64 `gorm:"column:buy_volume"`
	SellPrice       float64 `gorm:"column:sell_price"`
	SellVolume      int64   `gorm:"column:sell_volume"`
	VolumeRatio     float64 `gorm:"column:volume_ratio"`
}
