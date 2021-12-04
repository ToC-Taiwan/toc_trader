// Package global package global
package global

// SystemSwitch SystemSwitch
type SystemSwitch struct {
	Buy                          bool `json:"buy"`
	Sell                         bool `json:"sell"`
	SellFirst                    bool `json:"sell_first"`
	BuyLater                     bool `json:"buy_later"`
	UseBidAsk                    bool `json:"use_bid_ask"`
	MeanTimeTradeStockNum        int  `json:"mean_time_trade_stock_num"`
	MeanTimeReverseTradeStockNum int  `json:"mean_time_reverse_trade_stock_num"`
}

// TradeType TradeType
type TradeType int64

const (
	// ForwardTrade ForwardTrade
	ForwardTrade TradeType = iota + 1
	// ReverseTrade ReverseTrade
	ReverseTrade
)
