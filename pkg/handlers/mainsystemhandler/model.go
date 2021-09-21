// Package mainsystemhandler main handler
package mainsystemhandler

// UpdateTradeBotConditionBody UpdateTradeBotConditionBody
type UpdateTradeBotConditionBody struct {
	EnableBuy             bool `json:"enable_buy"`
	EnableSell            bool `json:"enable_sell"`
	UseBidAsk             bool `json:"use_bid_ask"`
	MeanTimeTradeStockNum int  `json:"mean_time_trade_stock_num"`
}
