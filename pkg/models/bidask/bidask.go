// Package bidask package bidask
package bidask

import (
	"sort"

	"gorm.io/gorm"
)

// BidAsk BidAsk
type BidAsk struct {
	gorm.Model  `json:"-" swaggerignore:"true"`
	StockNum    string  `gorm:"column:stock_num;index:idx_bidask"`
	TimeStamp   int64   `gorm:"column:timestamp;index:idx_bidask"`
	BidPrice1   float64 `gorm:"column:bid_price_1"`
	BidVolume1  int64   `gorm:"column:bid_volume_1"`
	DiffBidVol1 int64   `gorm:"column:diff_bid_vol_1"`
	BidPrice2   float64 `gorm:"column:bid_price_2"`
	BidVolume2  int64   `gorm:"column:bid_volume_2"`
	DiffBidVol2 int64   `gorm:"column:diff_bid_vol_2"`
	BidPrice3   float64 `gorm:"column:bid_price_3"`
	BidVolume3  int64   `gorm:"column:bid_volume_3"`
	DiffBidVol3 int64   `gorm:"column:diff_bid_vol_3"`
	BidPrice4   float64 `gorm:"column:bid_price_4"`
	BidVolume4  int64   `gorm:"column:bid_volume_4"`
	DiffBidVol4 int64   `gorm:"column:diff_bid_vol_4"`
	BidPrice5   float64 `gorm:"column:bid_price_5"`
	BidVolume5  int64   `gorm:"column:bid_volume_5"`
	DiffBidVol5 int64   `gorm:"column:diff_bid_vol_5"`
	AskPrice1   float64 `gorm:"column:ask_price_1"`
	AskVolume1  int64   `gorm:"column:ask_volume_1"`
	DiffAskVol1 int64   `gorm:"column:diff_ask_vol_1"`
	AskPrice2   float64 `gorm:"column:ask_price_2"`
	AskVolume2  int64   `gorm:"column:ask_volume_2"`
	DiffAskVol2 int64   `gorm:"column:diff_ask_vol_2"`
	AskPrice3   float64 `gorm:"column:ask_price_3"`
	AskVolume3  int64   `gorm:"column:ask_volume_3"`
	DiffAskVol3 int64   `gorm:"column:diff_ask_vol_3"`
	AskPrice4   float64 `gorm:"column:ask_price_4"`
	AskVolume4  int64   `gorm:"column:ask_volume_4"`
	DiffAskVol4 int64   `gorm:"column:diff_ask_vol_4"`
	AskPrice5   float64 `gorm:"column:ask_price_5"`
	AskVolume5  int64   `gorm:"column:ask_volume_5"`
	DiffAskVol5 int64   `gorm:"column:diff_ask_vol_5"`
	Suspend     int64   `gorm:"column:suspend"`
	Simtrade    int64   `gorm:"column:simtrade"`
}

// Tabler Tabler
type Tabler interface {
	TableName() string
}

// TableName TableName
func (BidAsk) TableName() string {
	return "tick_bidask"
}

// IsBestBid IsBestBid
func (c *BidAsk) IsBestBid() bool {
	tmpArr := append([]int64{}, c.BidVolume1, c.BidVolume2, c.BidVolume3, c.BidVolume4, c.BidVolume5)
	sort.Slice(tmpArr, func(i, j int) bool {
		return tmpArr[i] > tmpArr[j]
	})
	return tmpArr[0] == c.BidVolume1
}
