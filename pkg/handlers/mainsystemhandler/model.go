// Package mainsystemhandler main handler
package mainsystemhandler

// UpdateTradeBotConditionBody UpdateTradeBotConditionBody
type UpdateTradeBotConditionBody struct {
	EnableBuy             bool `json:"enable_buy"`
	EnableSell            bool `json:"enable_sell"`
	EnableSellFirst       bool `json:"enable_sell_first"`
	EnableBuyLater        bool `json:"enable_buy_later"`
	UseBidAsk             bool `json:"use_bid_ask"`
	MeanTimeTradeStockNum int  `json:"mean_time_trade_stock_num"`
}
