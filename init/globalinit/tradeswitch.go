// Package globalinit is init all global var
package globalinit

import (
	"os"

	"gitlab.tocraw.com/root/toc_trader/pkg/global"
)

func init() {
	global.TradeSwitch = global.SystemSwitch{
		Buy:                          true,
		Sell:                         true,
		SellFirst:                    false,
		BuyLater:                     true,
		UseBidAsk:                    false,
		MeanTimeTradeStockNum:        25,
		MeanTimeReverseTradeStockNum: 25,
	}

	deployment := os.Getenv("DEPLOYMENT")
	if deployment != "docker" {
		global.TradeSwitch.Buy = false
		global.TradeSwitch.SellFirst = false
	}
}
