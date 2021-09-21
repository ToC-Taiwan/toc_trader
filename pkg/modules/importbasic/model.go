// Package importbasic package importbasic
package importbasic

import (
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/stock"
)

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
