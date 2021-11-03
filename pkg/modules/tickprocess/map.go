// Package tickprocess package tickprocess
package tickprocess

import (
	"sync"

	"gitlab.tocraw.com/root/toc_trader/pkg/models/analyzeentiretick"
)

// AnalyzeEntireTickMap AnalyzeEntireTickMap
type AnalyzeEntireTickMap struct {
	tickMap map[string][]*analyzeentiretick.AnalyzeEntireTick
	mutex   sync.RWMutex
}

// SaveByStockNum SaveByStockNum
func (c *AnalyzeEntireTickMap) SaveByStockNum(stockNum string, tickArr []*analyzeentiretick.AnalyzeEntireTick) {
	if c.tickMap == nil {
		c.tickMap = make(map[string][]*analyzeentiretick.AnalyzeEntireTick)
	}
	c.mutex.Lock()
	c.tickMap[stockNum] = tickArr
	c.mutex.Unlock()
}

// GetAllTicks GetAllTicks
func (c *AnalyzeEntireTickMap) GetAllTicks() []*analyzeentiretick.AnalyzeEntireTick {
	var tmp []*analyzeentiretick.AnalyzeEntireTick
	if c.tickMap == nil {
		c.tickMap = make(map[string][]*analyzeentiretick.AnalyzeEntireTick)
	}
	c.mutex.RLock()
	for _, v := range c.tickMap {
		tmp = append(tmp, v...)
	}
	c.mutex.RUnlock()
	return tmp
}

// ClearAll ClearAll
func (c *AnalyzeEntireTickMap) ClearAll() {
	c.mutex.Lock()
	c.tickMap = make(map[string][]*analyzeentiretick.AnalyzeEntireTick)
	c.mutex.Unlock()
}
