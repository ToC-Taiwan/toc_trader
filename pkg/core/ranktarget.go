// Package core package core
package core

import (
	"time"

	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/choosetarget"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
)

func addRankTarget() {
	tick := time.NewTicker(15 * time.Second)
	for range tick.C {
		if !checkIsOpenTime() {
			continue
		}
		var count int
		if newTargetArr, err := choosetarget.GetTopTarget(10); err != nil {
			logger.Logger.Error(err)
			continue
		} else if time.Now().After(global.TradeDay.Add(1*time.Hour + 5*time.Minute)) {
			count = len(newTargetArr)
			if count != 0 {
				choosetarget.SubscribeTarget(newTargetArr)
				global.TargetArr = append(global.TargetArr, newTargetArr...)
			}
		}
		if count != 0 {
			logger.Logger.Infof("GetTopTarget %d", count)
		}
	}
}
