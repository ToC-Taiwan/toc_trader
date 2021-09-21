// Package subscribe package subscribe
package subscribe

import (
	"errors"
	"runtime/debug"

	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/pyresponse"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/streamtick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/bidaskprocess"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/streamtickprocess"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
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
			logger.Logger.Error(err.Error() + "\n" + string(debug.Stack()))
		}
	}()
	saveCh := make(chan []*streamtick.StreamTick, len(stockArr))
	go streamtickprocess.SaveStreamTicks(saveCh)

	for _, stockNum := range stockArr {
		ch := make(chan *streamtick.StreamTick)
		StreamTickChannelMap.Set(stockNum, ch)

		lastClose := global.StockCloseByDateMap.GetClose(stockNum, global.LastTradeDay.Format(global.ShortTimeLayout))
		go streamtickprocess.TickProcess(lastClose, ch, saveCh)
	}

	stocks := SubBody{
		StockNumArr: stockArr,
	}
	resp, err := global.RestyClient.R().
		SetBody(stocks).
		SetResult(&pyresponse.PyServerResponse{}).
		Post("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/subscribe/streamtick")
	if err != nil {
		panic(err)
	} else if resp.StatusCode() != 200 {
		panic("api fail")
	}
	res := *resp.Result().(*pyresponse.PyServerResponse)
	if res.Status != SucessStatus {
		panic("subscribe fail")
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
			logger.Logger.Error(err.Error() + "\n" + string(debug.Stack()))
		}
	}()

	for _, stock := range stockArr {
		go bidaskprocess.SaveBidAsk(stock)
	}
	stocks := SubBody{
		StockNumArr: stockArr,
	}
	resp, err := global.RestyClient.R().
		SetBody(stocks).
		SetResult(&pyresponse.PyServerResponse{}).
		Post("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/subscribe/bid-ask")
	if err != nil {
		panic(err)
	} else if resp.StatusCode() != 200 {
		panic("api fail")
	}
	res := *resp.Result().(*pyresponse.PyServerResponse)
	if res.Status != SucessStatus {
		panic("subscribe bidask fail")
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
			logger.Logger.Error(err.Error() + "\n" + string(debug.Stack()))
		}
	}()
	var url string
	if dataType == StreamType {
		url = "/pyapi/unsubscribeall/streamtick"
	} else if dataType == BidAsk {
		url = "/pyapi/unsubscribeall/bid-ask"
	}
	resp, err := global.RestyClient.R().
		SetResult(&pyresponse.PyServerResponse{}).
		Get("http://" + global.PyServerHost + ":" + global.PyServerPort + url)
	if err != nil {
		panic(err)
	} else if resp.StatusCode() != 200 {
		panic("api fail")
	}
	res := *resp.Result().(*pyresponse.PyServerResponse)
	if res.Status != SucessStatus {
		panic("unsubscribe fail")
	}
}
