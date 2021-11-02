// Package global all global var and struct
package global

import (
	"time"

	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
)

const (
	// ExitSignal ExitSignal
	ExitSignal int = 1
	// TradeYear TradeYear
	TradeYear int64 = 2021
	// OneTimeQuantity OneTimeQuantity
	OneTimeQuantity int64 = 1
	// LongTimeLayout LongTimeLayout
	LongTimeLayout string = "2006-01-02 15:04:05"
	// ShortTimeLayout ShortTimeLayout
	ShortTimeLayout string = "2006-01-02"
	// TradeInEndHour TradeInEndHour
	TradeInEndHour int = 10
	// TradeInEndMinute TradeInEndMinute
	TradeInEndMinute int = 15
	// TradeOutEndHour TradeOutEndHour
	TradeOutEndHour int = 13
	// TradeOutEndMinute TradeOutEndMinute
	TradeOutEndMinute int = 10
)

// ExitChannel ExitChannel
var ExitChannel chan int

// TradeDay TradeDay
var TradeDay time.Time

// TradeDayInEndTime TradeDayInEndTime
var TradeDayInEndTime time.Time

// TradeDayOutEndTime TradeDayOutEndTime
var TradeDayOutEndTime time.Time

// LastTradeDay LastTradeDay
var LastTradeDay time.Time

// LastLastTradeDay LastLastTradeDay
var LastLastTradeDay time.Time

// LastTradeDayArr LastTradeDayArr
var LastTradeDayArr []time.Time

// TargetArr TargetArr
var TargetArr []string

// HTTPPort HTTPPort
var HTTPPort string

// PyServerHost PyServerHost
var PyServerHost string

// PyServerPort PyServerPort
var PyServerPort string

// TradeSwitch TradeSwitch
var TradeSwitch SystemSwitch

// AllStockNameMap AllStockNameMap
var AllStockNameMap stringStringMutex

// StockCloseByDateMap StockCloseByDateMap
var StockCloseByDateMap stringStringFloat64Mutex

// BaseCond BaseCond
var BaseCond simulationcond.AnalyzeCondition

// ForwardCond ForwardCond
var ForwardCond simulationcond.AnalyzeCondition

// ReverseCond ReverseCond
var ReverseCond simulationcond.AnalyzeCondition
