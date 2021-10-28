// Package subscribe package subscribe
package subscribe

type subBody struct {
	StockNumArr []string `json:"stock_num_arr"`
}

// TickType TickType
type TickType int64

const (
	// StreamType StreamType
	StreamType TickType = iota + 1
	// BidAsk BidAsk
	BidAsk
)
