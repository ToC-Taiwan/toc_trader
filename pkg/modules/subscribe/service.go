// Package subscribe package subscribe
package subscribe

import (
	"errors"
	"runtime/debug"
	"sync"
	"time"

	"gitlab.tocraw.com/root/toc_trader/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/entiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/streamtick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/bidaskprocess"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/fetchentiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tickprocess"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tradebot"
	"gitlab.tocraw.com/root/toc_trader/pkg/sinopacapi"
)

// ForwardStreamTickChannelMap ForwardStreamTickChannelMap
var ForwardStreamTickChannelMap streamTickChannelMapMutexStruct

// ReverseStreamTickChannelMap ReverseStreamTickChannelMap
var ReverseStreamTickChannelMap streamTickChannelMapMutexStruct

// SimTradeChannel SimTradeChannel
var SimTradeChannel chan int

// SubStockStreamTick SubStockStreamTick
func SubStockStreamTick(stockArr []string) {
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
	saveCh := make(chan []*streamtick.StreamTick, len(stockArr)*3)
	go tickprocess.SaveStreamTicks(saveCh)
	for _, stockNum := range stockArr {
		lastClose := global.StockCloseByDateMap.GetClose(stockNum, global.LastTradeDay.Format(global.ShortTimeLayout))
		if lastClose == 0 {
			logger.GetLogger().Warnf("Stock %s has no lastClose", stockNum)
			continue
		}
		if global.TradeSwitch.Buy {
			forwardCh := make(chan *streamtick.StreamTick)
			ForwardStreamTickChannelMap.Set(stockNum, forwardCh)
			go tickprocess.ForwardTickProcess(lastClose, global.ForwardCond, forwardCh, saveCh)
		}
		if global.TradeSwitch.SellFirst {
			reverseCh := make(chan *streamtick.StreamTick)
			ReverseStreamTickChannelMap.Set(stockNum, reverseCh)
			go tickprocess.ReverseTickProcess(lastClose, global.ReverseCond, reverseCh)
		}
	}
	// fill missing ticks
	var wg sync.WaitGroup
	if tradebot.CheckIsOpenTime() {
		for _, v := range stockArr {
			wg.Add(1)
			go func(stock string) {
				defer wg.Done()
				var ticks []*entiretick.EntireTick
				if ticks, err = fetchentiretick.FetchByDate(stock, global.TradeDay.Format(global.ShortTimeLayout)); err != nil {
					logger.GetLogger().Error(err)
					return
				}
				forwardCh := *ForwardStreamTickChannelMap.GetChannelByStockNum(stock)
				reverseCh := *ReverseStreamTickChannelMap.GetChannelByStockNum(stock)
				for _, tick := range ticks {
					streamTick := tick.ToStreamTick()
					forwardCh <- streamTick
					reverseCh <- streamTick
				}
				logger.GetLogger().WithFields(map[string]interface{}{
					"StockNum": stock,
					"Count":    len(ticks),
				}).Info("Fill Missing Ticks Done")
			}(v)
		}
		wg.Wait()
	}
	for _, v := range stockArr {
		tickprocess.MissingTicksStatus.SetDone(v)
	}
	err = sinopacapi.GetAgent().SubStreamTick(stockArr)
	if err != nil {
		logger.GetLogger().Panic(err)
	}
}

// SubStockBidAsk SubStockBidAsk
func SubStockBidAsk(stockArr []string) {
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
	for _, stock := range stockArr {
		go bidaskprocess.SaveBidAsk(stock)
	}
	err = sinopacapi.GetAgent().SubBidAsk(stockArr)
	if err != nil {
		logger.GetLogger().Panic(err)
	}
}

// UnSubscribeStockAllByType UnSubscribeStockAllByType
func UnSubscribeStockAllByType(dataType sinopacapi.TickType) (err error) {
	err = sinopacapi.GetAgent().UnSubscribeAllByType(dataType)
	if err != nil {
		return err
	}
	return err
}

// SimTradeCollector SimTradeCollector
func SimTradeCollector() {
	SimTradeChannel = make(chan int)
	printMinute := time.Now().Minute()
	var count int
	for {
		simTrade := <-SimTradeChannel
		count += simTrade
		if time.Now().Minute() != printMinute {
			printMinute = time.Now().Minute()
			logger.GetLogger().WithFields(map[string]interface{}{
				"Count": count,
			}).Info("SimTrade")
		}
	}
}
