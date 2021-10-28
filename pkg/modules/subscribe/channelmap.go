// Package subscribe package subscribe
package subscribe

import (
	"sync"

	"gitlab.tocraw.com/root/toc_trader/pkg/models/streamtick"
)

type streamTickChannelMapMutexStruct struct {
	streamTickChannelMap map[string]*chan *streamtick.StreamTick
	mutex                sync.RWMutex
}

func (c *streamTickChannelMapMutexStruct) Set(stockNum string, ch chan *streamtick.StreamTick) {
	if c.streamTickChannelMap == nil {
		c.streamTickChannelMap = make(map[string]*chan *streamtick.StreamTick)
	}
	c.mutex.Lock()
	c.streamTickChannelMap[stockNum] = &ch
	c.mutex.Unlock()
}

func (c *streamTickChannelMapMutexStruct) GetChannelByStockNum(stockNum string) *chan *streamtick.StreamTick {
	var tmp *chan *streamtick.StreamTick
	c.mutex.RLock()
	tmp = c.streamTickChannelMap[stockNum]
	c.mutex.RUnlock()
	return tmp
}
