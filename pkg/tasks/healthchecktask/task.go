// Package healthchecktask package healthchecktask
package healthchecktask

import (
	"errors"
	"runtime/debug"
	"sync"

	"gitlab.tocraw.com/root/toc_trader/internal/healthcheck"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
)

var lock sync.Mutex

// Run Run
func Run() {
	lock.Lock()
	var err error
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
			logger.GetLogger().Error(err.Error() + "\n" + string(debug.Stack()))
		}
	}()
	defer lock.Unlock()
	if err = healthcheck.AskSinopacSRVRestart(); err != nil {
		logger.GetLogger().Panic(err)
	}
}
