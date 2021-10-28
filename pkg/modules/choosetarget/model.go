// Package choosetarget package choosetarget
package choosetarget

// FetchLastCountBody FetchLastCountBody
type FetchLastCountBody struct {
	StockNumArr []string `json:"stock_num_arr"`
}

// UpdateVolumeArrBody UpdateVolumeArrBody
type UpdateVolumeArrBody struct {
	StockNumArr []string `json:"stock_num_arr"`
}
