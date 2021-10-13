// Package fullrestart package fullrestart
package fullrestart

import (
	"sync"

	"gitlab.tocraw.com/root/toc_trader/tools/heartbeat"
)

var lock sync.RWMutex

// Run Run
func Run() {
	lock.Lock()
	defer lock.Unlock()
	if err := heartbeat.FullRestart(); err != nil {
		panic(err)
	}
}
