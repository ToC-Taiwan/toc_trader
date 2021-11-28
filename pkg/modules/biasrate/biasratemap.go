// Package biasrate package biasrate
package biasrate

import (
	"sync"
)

// MutexCache MutexCache
type MutexCache struct {
	rateMap map[string]map[string]float64
	mutex   sync.RWMutex
}

// Set Set
func (c *MutexCache) Set(stockNum, date string, biasRate float64) {
	if c.rateMap == nil {
		c.rateMap = make(map[string]map[string]float64)
	}
	if c.rateMap[stockNum] == nil {
		c.rateMap[stockNum] = make(map[string]float64)
	}
	c.mutex.Lock()
	c.rateMap[stockNum][date] = biasRate
	c.mutex.Unlock()
}

// GetBiasRate GetBiasRate
func (c *MutexCache) GetBiasRate(stockNum, date string) float64 {
	var tmp float64
	c.mutex.RLock()
	tmp = c.rateMap[stockNum][date]
	c.mutex.RUnlock()
	return tmp
}
