// Package sinopacsrv package sinopacsrv
package sinopacsrv

import (
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/stock"
)

const (
	// StatusSuccuss StatusSuccuss
	StatusSuccuss string = "success"
	// StatusFail StatusFail
	StatusFail string = "fail"
	// StatusNotFound StatusNotFound
	StatusNotFound string = "not found"
	// StatusAlready StatusAlready
	StatusAlready string = "already"
)

// OrderResponse OrderResponse
type OrderResponse struct {
	Status  string `json:"status"`
	OrderID string `json:"order_id"`
}

// StockLastCount StockLastCount
type StockLastCount struct {
	Date  string    `json:"date"`
	Code  string    `json:"code"`
	Time  []int64   `json:"time"`
	Close []float64 `json:"close"`
}

// SinoPacOrderStatus SinoPacOrderStatus
type SinoPacOrderStatus struct {
	Action    string  `json:"action"`
	Code      string  `json:"code"`
	ID        string  `json:"id"`
	Price     float64 `json:"price"`
	Quantity  int64   `json:"quantity"`
	Status    string  `json:"status"`
	OrderTime string  `json:"order_time"`
}

// SinoStatusResponse SinoStatusResponse
type SinoStatusResponse struct {
	Status string               `json:"status"`
	Data   []SinoPacOrderStatus `json:"data"`
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

// ToStock ToStock
func (c *FetchStockBody) ToStock() stock.Stock {
	global.AllStockNameMap.Set(c.Code, c.Name)
	var dayTradeBool bool
	if c.DayTrade == "Yes" {
		dayTradeBool = true
	} else {
		dayTradeBool = false
	}
	return stock.Stock{
		StockNum:  c.Code,
		StockName: c.Name,
		StockType: c.Exchange,
		DayTrade:  dayTradeBool,
		LastClose: c.Close,
		Category:  c.Category,
	}
}

// SinopacHealthStatus SinopacHealthStatus
type SinopacHealthStatus struct {
	Status      string `json:"status"`
	UpTimeMin   int64  `json:"up_time_min"`
	ServerToken string `json:"server_token"`
}
