// Package choosetarget package choosetarget
package choosetarget

// StockLastCount StockLastCount
type StockLastCount struct {
	Date  string    `json:"date"`
	Code  string    `json:"code"`
	Time  []int64   `json:"time"`
	Close []float64 `json:"close"`
}

// FetchLastCountBody FetchLastCountBody
type FetchLastCountBody struct {
	StockNumArr []string `json:"stock_num_arr"`
}

// UpdateVolumeArrBody UpdateVolumeArrBody
type UpdateVolumeArrBody struct {
	StockNumArr []string `json:"stock_num_arr"`
}

// SinoPacStockVolumeClose SinoPacStockVolumeClose
type SinoPacStockVolumeClose struct {
	Code   string  `json:"code"`
	Close  float64 `json:"close"`
	Volume int64   `json:"volume"`
}
