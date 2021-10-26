// Package core package core
package core

import (
	"os"
	"time"

	"github.com/manifoldco/promptui"
	"gitlab.tocraw.com/root/toc_trader/init/sysparminit"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/stock"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/targetstock"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/choosetarget"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/fetchentiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/importbasic"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/simulateprocess"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tradebot"
	"gitlab.tocraw.com/root/toc_trader/tools/db"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
)

// TradeProcess TradeProcess
func TradeProcess() {
	// Import all stock and update AllStockNameMap
	importbasic.ImportAllStock()
	// Development
	deployment := os.Getenv("DEPLOYMENT")
	if deployment != "docker" {
		// Send ip to sinopac srv
		sendCurrentIP()
		// Simulate
		prompt := promptui.Prompt{
			Label: "Simulate?(y/n)",
		}
		result, err := prompt.Run()
		if err != nil {
			panic(err)
		}
		if result == "y" {
			simulateprocess.Simulate()
		}
	}
	// Generate global target array and fetch entireTick
	var savedTarget []targetstock.Target
	if targets, err := choosetarget.GetTargetByVolumeRankByDate(global.LastTradeDay.Format(global.ShortTimeLayout), 200); err != nil {
		panic(err)
	} else {
		// UnSubscribeAll first
		choosetarget.UnSubscribeAll()
		// Subscribe all target
		global.TargetArr = targets
		choosetarget.SubscribeTarget(global.TargetArr)
		for i, v := range targets {
			logger.GetLogger().Infof("%s volume rank no. %d is %s\n", global.LastTradeDay.Format(global.ShortTimeLayout), i+1, global.AllStockNameMap.GetName(v))
		}
		tmp := []time.Time{global.LastTradeDay}
		fetchentiretick.FetchEntireTick(targets, tmp, global.CentralCond)
		if dbTarget, err := targetstock.GetTargetByTime(global.LastTradeDay, db.GetAgent()); err != nil {
			panic(err)
		} else if len(dbTarget) == 0 {
			logger.GetLogger().Info("Saving targets")
			targetStockArr, err := stock.GetStocksFromNumArr(targets, db.GetAgent())
			if err != nil {
				panic(err)
			}
			for _, v := range targetStockArr {
				savedTarget = append(savedTarget, targetstock.Target{
					LastTradeDay: global.LastTradeDay,
					Stock:        v,
				})
			}
			if err := targetstock.InsertMultiTarget(savedTarget, db.GetAgent()); err != nil {
				panic(err)
			}
		}
	}
	// Fetch Kbar
	kbarTradeDayArr, err := importbasic.GetLastNTradeDay(sysparminit.GlobalSettings.GetKbarPeriod())
	if err != nil {
		panic(err)
	}
	fetchentiretick.FetchKbar(global.TargetArr, kbarTradeDayArr[len(kbarTradeDayArr)-1], kbarTradeDayArr[0])
	logger.GetLogger().Info("FetchEntireTick and Kbar Done")

	// Check tradeday or target exist
	if len(global.LastTradeDayArr) == 0 || len(global.TargetArr) == 0 {
		logger.GetLogger().Warn("no trade day or no target")
	} else {
		// Background get trade record
		logger.GetLogger().Info("Background tasks starts")
		go tradebot.CheckOrderStatusLoop()
		// Monitor TSE001 Status
		go choosetarget.TSEProcess()
		// Add Top Rank Targets
		go addRankTarget()
	}
}
