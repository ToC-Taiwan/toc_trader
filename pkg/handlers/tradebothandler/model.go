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
