// Package tradebothandler tradebothandler
package tradebothandler

// ManualSellBody ManualSellBody
type ManualSellBody struct {
	StockNum string  `json:"stock_num"`
	Price    float64 `json:"price"`
}

// ManualBuyLaterBody ManualBuyLaterBody
type ManualBuyLaterBody struct {
	StockNum string  `json:"stock_num"`
	Price    float64 `json:"price"`
}

// TargetResponse TargetResponse
type TargetResponse struct {
	StockNum string  `json:"stock_num"`
	Close    float64 `json:"close"`
}

// UpdateTradeBotSwitchBody UpdateTradeBotSwitchBody
type UpdateTradeBotSwitchBody struct {
	EnableBuy                    bool `json:"enable_buy"`
	EnableSell                   bool `json:"enable_sell"`
	EnableSellFirst              bool `json:"enable_sell_first"`
	EnableBuyLater               bool `json:"enable_buy_later"`
	UseBidAsk                    bool `json:"use_bid_ask"`
	MeanTimeTradeStockNum        int  `json:"mean_time_trade_stock_num"`
	MeanTimeReverseTradeStockNum int  `json:"mean_time_reverse_trade_stock_num"`
}
