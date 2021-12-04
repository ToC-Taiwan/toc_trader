// Package sinopacapi package sinopacapi
package sinopacapi

// OrderBody OrderBody
type OrderBody struct {
	Stock    string  `json:"stock"`
	Price    float64 `json:"price"`
	Quantity int64   `json:"quantity"`
}

// OrderCancelBody OrderCancelBody
type OrderCancelBody struct {
	OrderID string `json:"order_id"`
}

// FetchLastCloseBody FetchLastCloseBody
type FetchLastCloseBody struct {
	StockNumArr []string `json:"stock_num_arr"`
	DateArr     []string `json:"date_arr"`
}

// FetchLastCountBody FetchLastCountBody
type FetchLastCountBody struct {
	StockNumArr []string `json:"stock_num_arr"`
}

// FetchKbarBody FetchKbarBody
type FetchKbarBody struct {
	StockNum  string `json:"stock_num"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// FetchBody FetchBody
type FetchBody struct {
	StockNum string `json:"stock_num"`
	Date     string `json:"date"`
}

// SubscribeBody SubscribeBody
type SubscribeBody struct {
	StockNumArr []string `json:"stock_num_arr"`
}

// FetchStockBody FetchStockBody
type FetchStockBody struct {
	Close    float64 `json:"close"`
	Code     string  `json:"code"`
	DayTrade string  `json:"day_trade"`
	Exchange string  `json:"exchange"`
	Name     string  `json:"name"`
	Updated  string  `json:"updated"`
	Category string  `json:"category"`
}

// LastCloseWithStockAndDate LastCloseWithStockAndDate
type LastCloseWithStockAndDate struct {
	StockNum string `json:"stock_num"`
	CloseArr []struct {
		Date  string  `json:"date"`
		Close float64 `json:"close"`
	} `json:"close_arr"`
}
