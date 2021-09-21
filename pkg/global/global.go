// Package global all global var and struct
package global

import (
	"time"

	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
)

const (
	// TradeYear TradeYear
	TradeYear int64 = 2021
	// OneTimeQuantity OneTimeQuantity
	OneTimeQuantity int64 = 1
	// LongTimeLayout LongTimeLayout
	LongTimeLayout string = "2006-01-02 15:04:05"
	// ShortTimeLayout ShortTimeLayout
	ShortTimeLayout string = "2006-01-02"
)

// GlobalDB GlobalDB
var GlobalDB *gorm.DB

// RestyClient RestyClient
var RestyClient *resty.Client

// ExitChannel ExitChannel
var ExitChannel chan string

// TradeDay TradeDay
var TradeDay time.Time

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

// AllStockNameMap AllStockNameMap
var AllStockNameMap stringStringMutex

// StockCloseByDateMap StockCloseByDateMap
var StockCloseByDateMap stringStringFloat64Mutex

// EnableBuy EnableBuy
var EnableBuy bool

// EnableSell EnableSell
var EnableSell bool

// UseBidAsk UseBidAsk
var UseBidAsk bool

// MeanTimeTradeStockNum MeanTimeTradeStockNum
var MeanTimeTradeStockNum int

// TickAnalyzeCondition TickAnalyzeCondition
var TickAnalyzeCondition AnalyzeCondition

// HistoryCloseCount HistoryCloseCount
var HistoryCloseCount int
