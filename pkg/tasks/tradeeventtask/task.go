// Package tradeeventtask is task for trade event
package tradeeventtask

import (
	"sync"

	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tradeeventprocess"
)

var lock sync.Mutex

// Run Run
func Run() {
	lock.Lock()
	defer lock.Unlock()
	if err := tradeeventprocess.CleanEvent(); err != nil {
		panic(err)
	}
}
