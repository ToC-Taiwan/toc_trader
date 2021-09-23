// Package tradeeventprocess is task for trade event
package tradeeventprocess

import (
	"errors"
	"runtime/debug"
	"sync"

	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tradeeventprocess"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
)

var lock sync.RWMutex

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
			logger.Logger.Error(err.Error() + "\n" + string(debug.Stack()))
		}
	}()
	defer lock.Unlock()
	if err := tradeeventprocess.CleanEvent(); err != nil {
		panic(err)
	}
}
