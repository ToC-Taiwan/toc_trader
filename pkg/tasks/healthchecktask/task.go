// Package healthchecktask package healthchecktask
package healthchecktask

import (
	"sync"

	"gitlab.tocraw.com/root/toc_trader/internal/healthcheck"
)

var lock sync.Mutex

// Run Run
func Run() {
	lock.Lock()
	defer lock.Unlock()
	if err := healthcheck.AskSinopacSRVRestart(); err != nil {
		panic(err)
	}
}
