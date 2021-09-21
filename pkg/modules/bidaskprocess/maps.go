// Package bidaskprocess package bidaskprocess
package bidaskprocess

import (
	"sync"

	"gitlab.tocraw.com/root/toc_trader/pkg/models/bidask"
)

// BidAskMutexStruct BidAskMutexStruct
type BidAskMutexStruct struct {
	BidAskMap map[string][]*bidask.BidAsk
	mutex     sync.RWMutex
}

// Append Append
func (c *BidAskMutexStruct) Append(record *bidask.BidAsk) {
	if c.BidAskMap == nil {
		c.BidAskMap = make(map[string][]*bidask.BidAsk)
	}
	c.mutex.Lock()
	c.BidAskMap[record.StockNum] = append(c.BidAskMap[record.StockNum], record)
	c.mutex.Unlock()
}

// KeepLastOne KeepLastOne
func (c *BidAskMutexStruct) KeepLastOne(stockNum string) {
	c.mutex.Lock()
	c.BidAskMap[stockNum] = c.BidAskMap[stockNum][len(c.BidAskMap[stockNum])-1:]
	c.mutex.Unlock()
}

// GetLastOneByStockNum GetLastOneByStockNum
func (c *BidAskMutexStruct) GetLastOneByStockNum(stockNum string) *bidask.BidAsk {
	var tmp *bidask.BidAsk
	if c.BidAskMap == nil {
		c.BidAskMap = make(map[string][]*bidask.BidAsk)
	}
	c.mutex.RLock()
	tmp = c.BidAskMap[stockNum][len(c.BidAskMap[stockNum])-1]
	c.mutex.RUnlock()
	return tmp
}

// GetArrByStockNum GetArrByStockNum
func (c *BidAskMutexStruct) GetArrByStockNum(stockNum string) []*bidask.BidAsk {
	var tmp []*bidask.BidAsk
	if c.BidAskMap == nil {
		c.BidAskMap = make(map[string][]*bidask.BidAsk)
	}
	c.mutex.RLock()
	tmp = c.BidAskMap[stockNum]
	c.mutex.RUnlock()
	return tmp
}

// GetCountByStockNum GetCountByStockNum
func (c *BidAskMutexStruct) GetCountByStockNum(stockNum string) int {
	var tmp int
	if c.BidAskMap == nil {
		c.BidAskMap = make(map[string][]*bidask.BidAsk)
	}
	c.mutex.RLock()
	tmp = len(c.BidAskMap[stockNum])
	c.mutex.RUnlock()
	return tmp
}
