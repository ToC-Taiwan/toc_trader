// Package fetchentiretick package fetchentiretick
package fetchentiretick

// FetchBody FetchBody
type FetchBody struct {
	StockNum string `json:"stock_num"`
	Date     string `json:"date"`
}

// FetchKbarBody FetchKbarBody
type FetchKbarBody struct {
	StockNum  string `json:"stock_num"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// FetchTSEKbarBody FetchTSEKbarBody
type FetchTSEKbarBody struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
