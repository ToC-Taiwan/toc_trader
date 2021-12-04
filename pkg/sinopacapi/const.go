// Package sinopacapi package sinopacapi
package sinopacapi

const (
	shortTimeLayout string = "2006-01-02"
)

const (
	// StatusSuccuss StatusSuccuss
	StatusSuccuss string = "success"
	// StatusFail StatusFail
	StatusFail string = "fail"
	// StatusHistoryNotFound StatusHistoryNotFound
	StatusHistoryNotFound string = "history not found"
	// StatusCancelOrderNotFound StatusCancelOrderNotFound
	StatusCancelOrderNotFound string = "cancel order not found"
	// StatusAlreadyCanceled StatusAlreadyCanceled
	StatusAlreadyCanceled string = "order already be canceled"
)

// OrderAction OrderAction
type OrderAction int64

const (
	// ActionBuy ActionBuy
	ActionBuy OrderAction = iota + 1
	// ActionSell ActionSell
	ActionSell
	// ActionSellFirst ActionSellFirst
	ActionSellFirst
)

// TickType TickType
type TickType int64

const (
	// StreamType StreamType
	StreamType TickType = iota + 1
	// BidAsk BidAsk
	BidAsk
)
