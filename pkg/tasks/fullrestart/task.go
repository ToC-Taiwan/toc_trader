// Package fullrestart package fullrestart
package fullrestart

import (
	"sync"

	"gitlab.tocraw.com/root/toc_trader/internal/healthcheck"
)

var lock sync.Mutex

// Run Run
func Run() {
	lock.Lock()
	defer lock.Unlock()
	if err := healthcheck.FullRestart(); err != nil {
		panic(err)
	}
}
