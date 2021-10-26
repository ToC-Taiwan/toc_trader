// Package subscribe package subscribe
package subscribe

import (
	"errors"
	"runtime/debug"

	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/streamtick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/bidaskprocess"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/streamtickprocess"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
	"gitlab.tocraw.com/root/toc_trader/tools/rest"
)

// SucessStatus SucessStatus
const SucessStatus string = "success"

// StreamTickChannelMap StreamTickChannelMap
var StreamTickChannelMap streamTickChannelMapMutexStruct

// SubStreamTick SubStreamTick
func SubStreamTick(stockArr []string) {
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
	go streamtickprocess.SaveStreamTicks(saveCh)

	for _, stockNum := range stockArr {
		lastClose := global.StockCloseByDateMap.GetClose(stockNum, global.LastTradeDay.Format(global.ShortTimeLayout))
		if lastClose == 0 {
			logger.GetLogger().Warnf("Stock %s has no lastClose", stockNum)
			continue
		}
		ch := make(chan *streamtick.StreamTick)
		StreamTickChannelMap.Set(stockNum, ch)
		go streamtickprocess.TickProcess(lastClose, global.CentralCond, ch, saveCh)
	}

	stocks := SubBody{
		StockNumArr: stockArr,
	}
	resp, err := rest.GetClient().R().
		SetBody(stocks).
		SetResult(&global.PyServerResponse{}).
		Post("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/subscribe/streamtick")
	if err != nil {
		panic(err)
	} else if resp.StatusCode() != 200 {
		panic("SubStreamTick api fail")
	}
	res := *resp.Result().(*global.PyServerResponse)
	if res.Status != SucessStatus {
		panic("Subscribe fail")
	}
}

// SubBidAsk SubBidAsk
func SubBidAsk(stockArr []string) {
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
	stocks := SubBody{
		StockNumArr: stockArr,
	}
	resp, err := rest.GetClient().R().
		SetBody(stocks).
		SetResult(&global.PyServerResponse{}).
		Post("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/subscribe/bid-ask")
	if err != nil {
		panic(err)
	} else if resp.StatusCode() != 200 {
		panic("SubBidAsk api fail")
	}
	res := *resp.Result().(*global.PyServerResponse)
	if res.Status != SucessStatus {
		panic("Subscribe bidask fail")
	}
}

// UnSubscribeAll UnSubscribeAll
func UnSubscribeAll(dataType TickType) {
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
	var url string
	switch {
	case dataType == StreamType:
		url = "/pyapi/unsubscribeall/streamtick"
	case dataType == BidAsk:
		url = "/pyapi/unsubscribeall/bid-ask"
	}
	resp, err := rest.GetClient().R().
		SetResult(&global.PyServerResponse{}).
		Get("http://" + global.PyServerHost + ":" + global.PyServerPort + url)
	if err != nil {
		panic(err)
	} else if resp.StatusCode() != 200 {
		panic("UnSubscribeAll api fail")
	}
	res := *resp.Result().(*global.PyServerResponse)
	if res.Status != SucessStatus {
		panic("Unsubscribe fail")
	}
}
