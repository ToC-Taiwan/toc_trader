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

func (c *stringStringMutex) GetValueByKey(key string) string {
	var tmp map[string]string
	c.mutex.RLock()
	tmp = c.sSMap
	c.mutex.RUnlock()
	if tmp == nil {
		tmp = make(map[string]string)
	}
	return tmp[key]
}
