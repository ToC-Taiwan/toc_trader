// Package streamtickprocess package streamtickprocess
package streamtickprocess

import (
	"sync"
)

// MissingTicksStatus MissingTicksStatus
var MissingTicksStatus statusMapMutexStruct

type statusMapMutexStruct struct {
	statusMap map[string]bool
	mutex     sync.RWMutex
}

func (c *statusMapMutexStruct) SetDone(stockNum string) {
	if c.statusMap == nil {
		c.statusMap = make(map[string]bool)
	}
	c.mutex.Lock()
	c.statusMap[stockNum] = true
	c.mutex.Unlock()
}

func (c *statusMapMutexStruct) CheckByStockNum(stockNum string) bool {
	var tmp bool
	c.mutex.RLock()
	tmp = c.statusMap[stockNum]
	c.mutex.RUnlock()
	return tmp
}
