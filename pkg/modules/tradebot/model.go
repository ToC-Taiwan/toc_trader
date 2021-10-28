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
	// SellFirstAction SellFirstAction
	SellFirstAction
)

type orderError string

const (
	// CancelFail CancelFail
	CancelFail orderError = "cancel fail"
	// CancelAlready CancelAlready
	CancelAlready orderError = "cancel already"
	// CancelNotFound CancelNotFound
	CancelNotFound orderError = "order not found"
)
