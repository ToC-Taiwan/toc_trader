// Package sinopacapi package sinopacapi
package sinopacapi

const (
	urlUpdateTraderIP    string = "/pyapi/system/tradebothost"
	urlFetchServerKey    string = "/pyapi/system/healthcheck"
	urlRestartSinopacSRV string = "/pyapi/system/restart"

	urlFetchAllStockDetail string = "/pyapi/basic/importstock"

	urlPlaceOrderBuy       string = "/pyapi/trade/buy"
	urlPlaceOrderSell      string = "/pyapi/trade/sell"
	urlPlaceOrderSellFirst string = "/pyapi/trade/sell_first"
	urlCancelOrder         string = "/pyapi/trade/cancel"

	urlFetchOrderStatus                 string = "/pyapi/trade/status"
	urlFetchAllSnapShot                 string = "/pyapi/basic/update/snapshot"
	urlFetchLastCloseByStockArrDateArr  string = "/pyapi/history/lastcount/multi-date"
	urlFetchStockCloseMapByStockDateArr string = "/pyapi/history/lastcount"
	urlFetchTSE001CloseByDate           string = "/pyapi/history/lastcount/tse"
	urlFetchVolumeRankByDate            string = "/pyapi/trade/volumerank"
	urlFetchKbarByDateRange             string = "/pyapi/history/kbar"
	urlFetchTSE001KbarByDate            string = "/pyapi/history/kbar/tse"
	urlFetchEntireTickByStockAndDate    string = "/pyapi/history/entiretick"

	urlSubStreamTick string = "/pyapi/subscribe/streamtick"
	urlSubBidAsk     string = "/pyapi/subscribe/bid-ask"

	urlUnSubscribeAllStream string = "/pyapi/unsubscribeall/streamtick"
	urlUnSubscribeAllBidAsk string = "/pyapi/unsubscribeall/bid-ask"
)
