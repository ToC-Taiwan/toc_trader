// Package tradebot package tradebot
package tradebot

import (
	"sync"
	"time"

	"gitlab.tocraw.com/root/toc_trader/pkg/models/traderecord"
)

type tradeRecordMutexMap struct {
	tMap  map[string]traderecord.TradeRecord
	mutex sync.RWMutex
}

// Set Set
func (c *tradeRecordMutexMap) Set(record traderecord.TradeRecord) {
	if c.tMap == nil {
		c.tMap = make(map[string]traderecord.TradeRecord)
	}
	c.mutex.Lock()
	c.tMap[record.StockNum] = record
	c.mutex.Unlock()
}

// Delete Delete
func (c *tradeRecordMutexMap) Delete(stockNum string) {
	if c.tMap == nil {
		c.tMap = make(map[string]traderecord.TradeRecord)
	}
	c.mutex.Lock()
	delete(c.tMap, stockNum)
	c.mutex.Unlock()
}

// ClearAll ClearAll
func (c *tradeRecordMutexMap) ClearAll() {
	c.mutex.Lock()
	c.tMap = make(map[string]traderecord.TradeRecord)
	c.mutex.Unlock()
}

// GetAll GetAll
func (c *tradeRecordMutexMap) GetAll() map[string]traderecord.TradeRecord {
	var tmp map[string]traderecord.TradeRecord
	if c.tMap == nil {
		c.tMap = make(map[string]traderecord.TradeRecord)
	}
	c.mutex.RLock()
	tmp = c.tMap
	c.mutex.RUnlock()
	return tmp
}

// GetCount GetCount
func (c *tradeRecordMutexMap) GetCount() int {
	var tmp int
	if c.tMap == nil {
		c.tMap = make(map[string]traderecord.TradeRecord)
	}
	c.mutex.RLock()
	tmp = len(c.tMap)
	c.mutex.RUnlock()
	return tmp
}

// GetOrderID GetOrderID
func (c *tradeRecordMutexMap) GetOrderID(stockNum string) string {
	var tmp string
	if c.tMap == nil {
		c.tMap = make(map[string]traderecord.TradeRecord)
	}
	c.mutex.RLock()
	tmp = c.tMap[stockNum].OrderID
	c.mutex.RUnlock()
	return tmp
}

// GetTradeTime GetTradeTime
func (c *tradeRecordMutexMap) GetTradeTime(stockNum string) time.Time {
	var tmp time.Time
	if c.tMap == nil {
		c.tMap = make(map[string]traderecord.TradeRecord)
	}
	c.mutex.RLock()
	tmp = c.tMap[stockNum].TradeTime
	c.mutex.RUnlock()
	return tmp
}

// CheckStockExist CheckStockExist
func (c *tradeRecordMutexMap) CheckStockExist(stockNum string) bool {
	var tmp bool
	if c.tMap == nil {
		c.tMap = make(map[string]traderecord.TradeRecord)
	}
	c.mutex.RLock()
	if _, ok := c.tMap[stockNum]; ok {
		tmp = true
	}
	c.mutex.RUnlock()
	return tmp
}

// GetClose GetClose
func (c *tradeRecordMutexMap) GetClose(stockNum string) float64 {
	var tmp float64
	if c.tMap == nil {
		c.tMap = make(map[string]traderecord.TradeRecord)
	}
	c.mutex.RLock()
	tmp = c.tMap[stockNum].Price
	c.mutex.RUnlock()
	return tmp
}

func (c *tradeRecordMutexMap) GetTotalBuyCost() int64 {
	if c.tMap == nil {
		c.tMap = make(map[string]traderecord.TradeRecord)
	}
	c.mutex.RLock()
	var cost int64
	for _, order := range c.tMap {
		cost += GetStockBuyCost(order.Price, order.Quantity)
	}
	c.mutex.RUnlock()
	return cost
}

func (c *tradeRecordMutexMap) GetTotalSellCost() int64 {
	if c.tMap == nil {
		c.tMap = make(map[string]traderecord.TradeRecord)
	}
	c.mutex.RLock()
	var cost int64
	for _, order := range c.tMap {
		cost += GetStockSellCost(order.Price, order.Quantity)
	}
	c.mutex.RUnlock()
	return cost
}

func (c *tradeRecordMutexMap) GetTotalCostBack() int64 {
	if c.tMap == nil {
		c.tMap = make(map[string]traderecord.TradeRecord)
	}
	c.mutex.RLock()
	var cost int64
	for _, order := range c.tMap {
		cost += GetStockTradeFeeDiscount(order.Price, order.Quantity)
	}
	c.mutex.RUnlock()
	return cost
}