// Package targethandler package targethandler
package targethandler

// TargetResponse TargetResponse
type TargetResponse struct {
	StockNum string  `json:"stock_num"`
	Close    float64 `json:"close"`
}
