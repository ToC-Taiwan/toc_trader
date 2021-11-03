// Package subscribe package subscribe
package subscribe

import (
	"errors"
	"net/http"
	"runtime/debug"

	"gitlab.tocraw.com/root/toc_trader/external/sinopacsrv"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/internal/restful"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/entiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/streamtick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/bidaskprocess"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/fetchentiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/streamtickprocess"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tradebot"
)

// ForwardStreamTickChannelMap ForwardStreamTickChannelMap
var ForwardStreamTickChannelMap streamTickChannelMapMutexStruct

// ReverseStreamTickChannelMap ReverseStreamTickChannelMap
var ReverseStreamTickChannelMap streamTickChannelMapMutexStruct

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

		forwardCh := make(chan *streamtick.StreamTick)
		ForwardStreamTickChannelMap.Set(stockNum, forwardCh)
		go streamtickprocess.ForwardTickProcess(lastClose, global.ForwardCond, forwardCh, saveCh)

		reverseCh := make(chan *streamtick.StreamTick)
		ReverseStreamTickChannelMap.Set(stockNum, reverseCh)
		go streamtickprocess.ReverseTickProcess(lastClose, global.ReverseCond, reverseCh)
	}
	// fill missing ticks
	if tradebot.CheckIsOpenTime() {
		for _, stock := range stockArr {
			var ticks []*entiretick.EntireTick
			if ticks, err = fetchentiretick.FetchByDate(stock, global.TradeDay.Format(global.ShortTimeLayout)); err != nil {
				logger.GetLogger().Error(err)
				return
			}
			for _, tick := range ticks {
				*ForwardStreamTickChannelMap.GetChannelByStockNum(stock) <- tick.ToStreamTick()
				*ReverseStreamTickChannelMap.GetChannelByStockNum(stock) <- tick.ToStreamTick()
			}
			logger.GetLogger().WithFields(map[string]interface{}{
				"StockNum": stock,
				"Count":    len(ticks),
			}).Info("Fill Missing Ticks Done")
		}
	}
	for _, v := range stockArr {
		streamtickprocess.MissingTicksStatus.SetDone(v)
	}
	stocks := subBody{
		StockNumArr: stockArr,
	}
	resp, err := restful.GetClient().R().
		SetBody(stocks).
		SetResult(&sinopacsrv.OrderResponse{}).
		Post("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/subscribe/streamtick")
	if err != nil {
		panic(err)
	} else if resp.StatusCode() != http.StatusOK {
		panic("SubStreamTick api fail")
	}
	res := *resp.Result().(*sinopacsrv.OrderResponse)
	if res.Status != sinopacsrv.StatusSuccuss {
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
	stocks := subBody{
		StockNumArr: stockArr,
	}
	resp, err := restful.GetClient().R().
		SetBody(stocks).
		SetResult(&sinopacsrv.OrderResponse{}).
		Post("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/subscribe/bid-ask")
	if err != nil {
		panic(err)
	} else if resp.StatusCode() != http.StatusOK {
		panic("SubBidAsk api fail")
	}
	res := *resp.Result().(*sinopacsrv.OrderResponse)
	if res.Status != sinopacsrv.StatusSuccuss {
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
	resp, err := restful.GetClient().R().
		SetResult(&sinopacsrv.OrderResponse{}).
		Get("http://" + global.PyServerHost + ":" + global.PyServerPort + url)
	if err != nil {
		panic(err)
	} else if resp.StatusCode() != http.StatusOK {
		panic("UnSubscribeAll api fail")
	}
	res := *resp.Result().(*sinopacsrv.OrderResponse)
	if res.Status != sinopacsrv.StatusSuccuss {
		panic("Unsubscribe fail")
	}
}
