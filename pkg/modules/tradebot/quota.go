// Package tradebot package tradebot
package tradebot

import (
	"math"

	"gitlab.tocraw.com/root/toc_trader/internal/db"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/traderecord"
)

// TradeQuota TradeQuota
var TradeQuota int64 = 1000000

const (
	// TradeFeeRatio TradeFeeRatio
	TradeFeeRatio float64 = 0.001425
	// FeeDiscount FeeDiscount
	FeeDiscount float64 = 0.35
	// TradeTaxRatio TradeTaxRatio
	TradeTaxRatio float64 = 0.0015
)

// InitStartUpQuota InitStartUpQuota
func InitStartUpQuota() (err error) {
	realOrder, err := traderecord.GetAllorderByDayTime(global.TradeDay, db.GetAgent())
	if err != nil {
		return err
	}
	for _, v := range realOrder {
		if v.Action == 1 && v.Status == 6 {
			TradeQuota -= GetStockBuyCost(v.Price, v.Quantity)
		}
	}
	return err
}

// GetStockBuyCost GetStockBuyCost
func GetStockBuyCost(price float64, qty int64) int64 {
	return int64(math.Ceil(price*float64(qty)*1000) + math.Floor(price*float64(qty)*1000*TradeFeeRatio))
}

// GetStockSellCost GetStockSellCost
func GetStockSellCost(price float64, qty int64) int64 {
	return int64(math.Ceil(price*float64(qty)*1000) - math.Floor(price*float64(qty)*1000*TradeFeeRatio) - math.Floor(price*float64(qty)*1000*TradeTaxRatio))
}

// GetStockTradeFeeDiscount GetStockTradeFeeDiscount
func GetStockTradeFeeDiscount(price float64, qty int64) int64 {
	return int64(math.Floor(price*float64(qty)*1000*TradeFeeRatio) * (1 - FeeDiscount))
}
