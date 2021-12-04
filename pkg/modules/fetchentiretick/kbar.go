// Package fetchentiretick package fetchentiretick
package fetchentiretick

import (
	"errors"
	"runtime/debug"
	"sync"
	"time"

	"gitlab.tocraw.com/root/toc_trader/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/database"
	"gitlab.tocraw.com/root/toc_trader/pkg/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/kbar"
	"gitlab.tocraw.com/root/toc_trader/pkg/sinopacapi"
)

// GetAndSaveKbar GetAndSaveKbar
func GetAndSaveKbar(stockNumArr []string, start, end time.Time) {
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
	saveCh := make(chan *kbar.Kbar)
	go kbarSaver(saveCh)
	var wg sync.WaitGroup
	for _, v := range stockNumArr {
		wg.Add(1)
		stock := v
		go func(stockNum string) {
			defer wg.Done()
			var exist bool
			exist, err = kbar.CheckExistByStockAndDateRange(stockNum, start, end, database.GetAgent())
			if err != nil {
				logger.GetLogger().Panic(err)
			}
			if !exist {
				if err = kbar.DeleteByStockNum(stockNum, database.GetAgent()); err != nil {
					logger.GetLogger().Panic(err)
				}
				logger.GetLogger().WithFields(map[string]interface{}{
					"StockNum": stockNum,
					"From":     start.Format(global.ShortTimeLayout),
					"To":       end.Format(global.ShortTimeLayout),
				}).Info("Fetching Kbar")
				if data, err := sinopacapi.GetAgent().FetchKbarByDateRange(stockNum, start, end); err != nil {
					logger.GetLogger().Panic(err)
				} else {
					for _, v := range data {
						tmp := &kbar.Kbar{
							StockNum:  stockNum,
							TimeStamp: v.GetTs(),
							Close:     v.GetClose(),
							Open:      v.GetOpen(),
							High:      v.GetHigh(),
							Low:       v.GetLow(),
							Volume:    v.GetVolume(),
						}
						saveCh <- tmp
					}
				}
			} else {
				logger.GetLogger().WithFields(map[string]interface{}{
					"StockNum": stockNum,
				}).Info("Kbar Already Exist")
			}
		}(stock)
	}
	wg.Wait()
	close(saveCh)
}

func kbarSaver(saveCh chan *kbar.Kbar) {
	var tmp []*kbar.Kbar
	for {
		kbarData, ok := <-saveCh
		if !ok {
			if err := kbar.InsertMultiRecord(tmp, database.GetAgent()); err != nil {
				logger.GetLogger().Error(err)
			}
			return
		}
		tmp = append(tmp, kbarData)
		if len(tmp) >= 2000 {
			if err := kbar.InsertMultiRecord(tmp, database.GetAgent()); err != nil {
				logger.GetLogger().Error(err)
				continue
			}
			tmp = []*kbar.Kbar{}
		}
	}
}
