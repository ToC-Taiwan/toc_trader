// Package simulateprocess package simulateprocess
package simulateprocess

import (
	"sync"
)

type targetArrMutex struct {
	arrMap map[string][]string
	mutex  sync.RWMutex
}

// saveByStockNum saveByStockNum
func (c *targetArrMutex) saveByDate(date string, targetArr []string) {
	if c.arrMap == nil {
		c.arrMap = make(map[string][]string)
	}
	c.mutex.Lock()
	c.arrMap[date] = targetArr
	c.mutex.Unlock()
}

// saveByStockNum saveByStockNum
func (c *targetArrMutex) getArrByDate(date string) (targetArr []string) {
	var tmp []string
	if c.arrMap == nil {
		c.arrMap = make(map[string][]string)
	}
	c.mutex.RLock()
	tmp = c.arrMap[date]
	c.mutex.RUnlock()
	return tmp
}

// clearAll clearAll
func (c *targetArrMutex) clearAll() {
	c.mutex.RLock()
	c.arrMap = make(map[string][]string)
	c.mutex.RUnlock()
}
