// Package simulateprocess package simulateprocess
package simulateprocess

import (
	"sync"

	"gitlab.tocraw.com/root/toc_trader/pkg/models/entiretick"
)

// entireTickMap entireTickMap
type entireTickMap struct {
	tickMap map[string]map[string][]*entiretick.EntireTick
	mutex   sync.RWMutex
}

// saveByStockNum saveByStockNum
func (c *entireTickMap) saveByStockNumAndDate(stockNum, date string, tickArr []*entiretick.EntireTick) {
	if c.tickMap == nil {
		c.tickMap = make(map[string]map[string][]*entiretick.EntireTick)
	}
	if c.tickMap[stockNum] == nil {
		c.tickMap[stockNum] = make(map[string][]*entiretick.EntireTick)
	}
	c.mutex.Lock()
	c.tickMap[stockNum][date] = tickArr
	c.mutex.Unlock()
}

func (c *entireTickMap) getAllTicksByStockNumAndDate(stockNum, date string) []*entiretick.EntireTick {
	var tmp []*entiretick.EntireTick
	if c.tickMap == nil {
		c.tickMap = make(map[string]map[string][]*entiretick.EntireTick)
	}
	if c.tickMap[stockNum] == nil {
		c.tickMap[stockNum] = make(map[string][]*entiretick.EntireTick)
	}
	c.mutex.RLock()
	tmp = c.tickMap[stockNum][date]
	c.mutex.RUnlock()
	return tmp
}

// clearAll clearAll
func (c *entireTickMap) clearAll() {
	c.mutex.Lock()
	c.tickMap = make(map[string]map[string][]*entiretick.EntireTick)
	c.mutex.Unlock()
}
