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
		SellFirst:                    true,
		BuyLater:                     true,
		UseBidAsk:                    false,
		MeanTimeTradeStockNum:        3,
		MeanTimeReverseTradeStockNum: 3,
	}

	deployment := os.Getenv("DEPLOYMENT")
	if deployment != "docker" {
		global.TradeSwitch.Buy = false
		global.TradeSwitch.SellFirst = false
	}
}
