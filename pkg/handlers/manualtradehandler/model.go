// Package manualtradehandler package manualtradehandler
package manualtradehandler

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
