// Package sinopacapi package sinopacapi
package sinopacapi

// ResponseHealthStatus ResponseHealthStatus
type ResponseHealthStatus struct {
	Status      string `json:"status"`
	UpTimeMin   int64  `json:"up_time_min"`
	ServerToken string `json:"server_token"`
}

// OrderResponse OrderResponse
type OrderResponse struct {
	Status  string `json:"status"`
	OrderID string `json:"order_id"`
}

// ResponseCommon ResponseCommon
type ResponseCommon struct {
	Status string `json:"status"`
}

// StockLastCount StockLastCount
type StockLastCount struct {
	Date  string    `json:"date"`
	Code  string    `json:"code"`
	Time  []int64   `json:"time"`
	Close []float64 `json:"close"`
}
