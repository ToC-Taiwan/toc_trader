// Package fullrestart package fullrestart
package fullrestart

import (
	"sync"

	"gitlab.tocraw.com/root/toc_trader/pkg/core"
)

var lock sync.RWMutex

// Run Run
func Run() {
	lock.Lock()
	defer lock.Unlock()
	if err := core.FullRestart(); err != nil {
		panic(err)
	}
}
