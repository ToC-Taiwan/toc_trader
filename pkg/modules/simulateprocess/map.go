// Package simulateprocess package simulateprocess
package simulateprocess

import (
	"sync"

	"gitlab.tocraw.com/root/toc_trader/pkg/models/entiretick"
)

// entireTickMap entireTickMap
type entireTickMap struct {
	tickMap map[string][]*entiretick.EntireTick
	mutex   sync.RWMutex
}

// saveByStockNum saveByStockNum
func (c *entireTickMap) saveByStockNum(stockNum string, tickArr []*entiretick.EntireTick) {
	if c.tickMap == nil {
		c.tickMap = make(map[string][]*entiretick.EntireTick)
	}
	c.mutex.Lock()
	c.tickMap[stockNum] = tickArr
	c.mutex.Unlock()
}

func (c *entireTickMap) getAllTicksByStockNum(stockNum string) []*entiretick.EntireTick {
	var tmp []*entiretick.EntireTick
	if c.tickMap == nil {
		c.tickMap = make(map[string][]*entiretick.EntireTick)
	}
	c.mutex.Lock()
	tmp = c.tickMap[stockNum]
	c.mutex.Unlock()
	return tmp
}

// saveByStockNum saveByStockNum
func (c *entireTickMap) getAllTicksMap() map[string][]*entiretick.EntireTick {
	var tmp map[string][]*entiretick.EntireTick
	if c.tickMap == nil {
		c.tickMap = make(map[string][]*entiretick.EntireTick)
	}
	c.mutex.Lock()
	tmp = c.tickMap
	c.mutex.Unlock()
	return tmp
}
