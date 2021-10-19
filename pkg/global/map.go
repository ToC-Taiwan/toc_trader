// Package global all global var and struct
package global

import "sync"

type stringStringMutex struct {
	sSMap map[string]string
	mutex sync.RWMutex
}

func (c *stringStringMutex) Set(key, value string) {
	if c.sSMap == nil {
		c.sSMap = make(map[string]string)
	}
	c.mutex.Lock()
	c.sSMap[key] = value
	c.mutex.Unlock()
}

func (c *stringStringMutex) Get() map[string]string {
	var tmp map[string]string
	c.mutex.RLock()
	tmp = c.sSMap
	c.mutex.RUnlock()
	if tmp == nil {
		tmp = make(map[string]string)
	}
	return tmp
}

func (c *stringStringMutex) GetName(key string) string {
	var tmp map[string]string
	c.mutex.RLock()
	tmp = c.sSMap
	c.mutex.RUnlock()
	if tmp == nil {
		tmp = make(map[string]string)
	}
	return tmp[key]
}

type stringStringFloat64Mutex struct {
	ssFMap map[string]map[string]float64
	mutex  sync.RWMutex
}

func (c *stringStringFloat64Mutex) Set(date, stockNum string, value float64) {
	if c.ssFMap == nil {
		c.ssFMap = make(map[string]map[string]float64)
	}
	if c.ssFMap[date] == nil {
		c.ssFMap[date] = make(map[string]float64)
	}
	c.mutex.Lock()
	c.ssFMap[date][stockNum] = value
	c.mutex.Unlock()
}

func (c *stringStringFloat64Mutex) Get() map[string]map[string]float64 {
	var tmp map[string]map[string]float64
	c.mutex.RLock()
	tmp = c.ssFMap
	c.mutex.RUnlock()
	if tmp == nil {
		tmp = make(map[string]map[string]float64)
	}
	return tmp
}

func (c *stringStringFloat64Mutex) GetClose(stockNum, date string) float64 {
	var tmp float64
	c.mutex.RLock()
	tmp = c.ssFMap[date][stockNum]
	c.mutex.RUnlock()
	return tmp
}
