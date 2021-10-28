// Package traderecord package traderecord
package traderecord

import (
	"time"

	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gorm.io/gorm"
)

// TradeRecord TradeRecord
type TradeRecord struct {
	gorm.Model
	StockNum  string    `gorm:"column:stock_num;index:idx_traderecord"`
	StockName string    `gorm:"column:stock_name"`
	Action    int64     `gorm:"column:action"`
	Price     float64   `gorm:"column:price"`
	Quantity  int64     `gorm:"column:quantity"`
	Status    int64     `gorm:"column:status"`
	OrderID   string    `gorm:"column:order_id;index:idx_traderecord"`
	OrderTime time.Time `gorm:"column:order_time"`
	BuyCost   int64     `gorm:"-"`
	TradeTime time.Time `gorm:"-"`
}

// Tabler Tabler
type Tabler interface {
	TableName() string
}

// TableName TableName
func (TradeRecord) TableName() string {
	return "trade_record"
}

// ActionListMap ActionListMap
var ActionListMap = map[string]int64{
	"Buy":  1,
	"Sell": 2,
}

// StatusListMap StatusListMap
var StatusListMap = map[string]int64{
	"PendingSubmit": 1, // 傳送中
	"PreSubmitted":  2, // 預約單
	"Submitted":     3, // 傳送成功
	"Failed":        4, // 失敗
	"Canceled":      5, // 已刪除
	"Filled":        6, // 完全成交
	"Filling":       7, // 部分成交
}

const (
	// PendingSubmit PendingSubmit
	PendingSubmit string = "PendingSubmit"
	// PreSubmitted PreSubmitted
	PreSubmitted string = "PreSubmitted"
	// Submitted Submitted
	Submitted string = "Submitted"
	// Failed Failed
	Failed string = "Failed"
	// Canceled Canceled
	Canceled string = "Canceled"
	// Filled Filled
	Filled string = "Filled"
	// Filling Filling
	Filling string = "Filling"
)

// ToTradeRecordFromProto ToTradeRecordFromProto
func (c *TradeRecordArrProto) ToTradeRecordFromProto() (record []*TradeRecord, err error) {
	for _, v := range c.Data {
		var orderTime time.Time
		name := global.AllStockNameMap.GetName(v.Code)
		orderTime, err = time.ParseInLocation(global.LongTimeLayout, v.OrderTime, time.Local)
		if err != nil {
			return record, err
		}
		status := StatusListMap[v.Status]
		action := ActionListMap[v.Action]
		tmp := TradeRecord{
			StockNum:  v.Code,
			StockName: name,
			Action:    action,
			Price:     v.Price,
			Quantity:  v.Quantity,
			Status:    status,
			OrderID:   v.Id,
			OrderTime: orderTime,
		}
		record = append(record, &tmp)
	}
	return record, err
}
