// Package healthchecktask package healthchecktask
package healthchecktask

import (
	"errors"
	"runtime/debug"
	"sync"

	"gitlab.tocraw.com/root/toc_trader/pkg/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/sinopacapi"
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
	if err = sinopacapi.GetAgent().RestartSinopacSRV(); err != nil {
		logger.GetLogger().Panic(err)
	}
}
