// Package tradeeventtask is task for trade event
package tradeeventtask

import (
	"errors"
	"runtime/debug"
	"sync"

	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tradeeventprocess"
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
	if err := tradeeventprocess.CleanEvent(); err != nil {
		panic(err)
	}
}
