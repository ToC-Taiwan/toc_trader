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
	"gitlab.tocraw.com/root/toc_trader/pkg/models/entiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/importbasic"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tickprocess"
	"gitlab.tocraw.com/root/toc_trader/pkg/sinopacapi"
)

var wg sync.WaitGroup

// FetchEntireTick FetchEntireTick
func FetchEntireTick(stockNumArr []string, dateArr []time.Time, cond simulationcond.AnalyzeCondition) {
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
	saveCh := make(chan []*entiretick.EntireTick, len(stockNumArr))
	go tickprocess.SaveEntireTicks(saveCh)
	for _, d := range dateArr {
		for _, s := range stockNumArr {
			var rows int64
			rows, err = entiretick.GetCntByStockAndDate(s, d.Format(global.ShortTimeLayout), database.GetAgent())
			if err != nil {
				logger.GetLogger().Panic(err)
			} else {
				if rows > 0 {
					logger.GetLogger().WithFields(map[string]interface{}{
						"Stock": s,
						"Date":  d.Format(global.ShortTimeLayout),
					}).Info("EntireTick Already Exist")
					continue
				} else {
					wg.Add(1)
					go GetEntireTickByStockAndDate(s, d.Format(global.ShortTimeLayout), cond, saveCh)
				}
			}
		}
		wg.Wait()
	}
	close(saveCh)
}

// GetEntireTickByStockAndDate GetEntireTickByStockAndDate
func GetEntireTickByStockAndDate(stockNum, date string, cond simulationcond.AnalyzeCondition, saveCh chan []*entiretick.EntireTick) {
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
	logger.GetLogger().WithFields(map[string]interface{}{
		"StockNum": stockNum,
		"Date":     date,
	}).Info("Fetching Entiretick")
	var res []*sinopacapi.EntireTickProto
	res, err = sinopacapi.GetAgent().FetchEntireTickByStockAndDate(stockNum, date)
	if err != nil {
		logger.GetLogger().Panic(err)
	}
	ch := make(chan *entiretick.EntireTick, len(res))
	var lastTradeDay time.Time
	lastTradeDay, err = importbasic.GetLastTradeDayByDate(date)
	if err != nil {
		logger.GetLogger().Panic(err)
	}
	var simulateMap tickprocess.AnalyzeEntireTickMap
	lastClose := global.StockCloseByDateMap.GetClose(stockNum, lastTradeDay.Format(global.ShortTimeLayout))
	if lastClose != 0 {
		go tickprocess.TickProcess(stockNum, lastClose, cond, ch, &wg, saveCh, false, &simulateMap)
	} else {
		logger.GetLogger().Warnf("%s has no %s's close", stockNum, date)
	}
	for _, v := range res {
		tick := &entiretick.EntireTick{
			StockNum:  stockNum,
			Close:     v.GetClose(),
			TickType:  v.GetTickType(),
			Volume:    v.GetVolume(),
			BidPrice:  v.GetBidPrice(),
			BidVolume: v.GetBidVolume(),
			AskPrice:  v.GetAskPrice(),
			AskVolume: v.GetAskVolume(),
			TimeStamp: v.GetTs(),
			Open:      0,
			High:      0,
			Low:       0,
		}
		ch <- tick
	}
	close(ch)
}

// FetchByDate FetchByDate
func FetchByDate(stockNum, date string) (data []*entiretick.EntireTick, err error) {
	var res []*sinopacapi.EntireTickProto
	res, err = sinopacapi.GetAgent().FetchEntireTickByStockAndDate(stockNum, date)
	if err != nil {
		return data, err
	}
	for _, v := range res {
		tick := &entiretick.EntireTick{
			StockNum:  stockNum,
			Close:     v.GetClose(),
			TickType:  v.GetTickType(),
			Volume:    v.GetVolume(),
			BidPrice:  v.GetBidPrice(),
			BidVolume: v.GetBidVolume(),
			AskPrice:  v.GetAskPrice(),
			AskVolume: v.GetAskVolume(),
			TimeStamp: v.GetTs(),
			Open:      0,
			High:      0,
			Low:       0,
		}
		data = append(data, tick)
	}
	return data, err
}
