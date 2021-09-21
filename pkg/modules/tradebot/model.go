// Package tradebot package tradebot
package tradebot

// OrderBody OrderBody
type OrderBody struct {
	Stock    string  `json:"stock"`
	Price    float64 `json:"price"`
	Quantity int64   `json:"quantity"`
}

// CancelBody CancelBody
type CancelBody struct {
	OrderID string `json:"order_id"`
}

// OrderAction OrderAction
type OrderAction int64

const (
	// BuyAction BuyAction
	BuyAction OrderAction = iota + 1
	// SellAction SellAction
	SellAction
)
