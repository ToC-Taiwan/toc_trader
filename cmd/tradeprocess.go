// Package main package main
package main

import (
	"os"
	"time"

	"github.com/manifoldco/promptui"
	"gitlab.tocraw.com/root/toc_trader/global"
	"gitlab.tocraw.com/root/toc_trader/init/sysparminit"
	"gitlab.tocraw.com/root/toc_trader/pkg/database"
	"gitlab.tocraw.com/root/toc_trader/pkg/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulate"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/stock"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/targetstock"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/choosetarget"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/fetchentiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/importbasic"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/simulateprocess"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/subscribe"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tickprocess"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tradebot"
)

// TradeProcess TradeProcess
func TradeProcess() {
	var err error
	// Check tradeday or target exist
	if len(global.LastTradeDayArr) == 0 {
		logger.GetLogger().Panic("No trade day")
	}
	// Import all stock and update AllStockNameMap
	if err = importbasic.ImportAllStock(); err != nil {
		logger.GetLogger().Panic(err)
	}
	// Clear stram tick in database
	if err = tickprocess.DeleteAllStreamTicks(); err != nil {
		logger.GetLogger().Panic(err)
	}
	// Chekc if is in debug mode and simulation
	simulatationEntry()
	// check db targets and save if no record in db
	dbTarget, err := targetstock.GetTargetByTime(global.TradeDay, database.GetAgent())
	if err != nil {
		logger.GetLogger().Panic(err)
	}
	if len(dbTarget) != 0 {
		logger.GetLogger().Info("Targets Already Exist")
		for _, v := range dbTarget {
			global.TargetArr = append(global.TargetArr, v.Stock.StockNum)
		}
	} else {
		logger.GetLogger().Info("Get Targets by Volume Rank")
		var targets []string
		targets, err = choosetarget.GetVolumeRankByDate(global.LastTradeDay.Format(global.ShortTimeLayout), 200)
		if err != nil {
			logger.GetLogger().Panic(err)
		}
		global.TargetArr = targets
		var savedTarget []targetstock.Target
		var targetStockArr []stock.Stock
		targetStockArr, err = stock.GetStocksFromNumArr(targets, database.GetAgent())
		if err != nil {
			logger.GetLogger().Panic(err)
		}
		for i, v := range targetStockArr {
			savedTarget = append(savedTarget, targetstock.Target{
				TradeDay: global.TradeDay,
				Stock:    v,
				Rank:     int64(i) + 1,
			})
		}
		if err = targetstock.InsertMultiRecord(savedTarget, database.GetAgent()); err != nil {
			logger.GetLogger().Panic(err)
		}
	}
	if len(global.TargetArr) == 0 {
		logger.GetLogger().Panic("No target")
	}
	// UnSubscribeAll first
	if err = choosetarget.UnSubscribeAllFromSinopac(); err != nil {
		logger.GetLogger().Panic(err)
	}
	// Simtrade collect
	go subscribe.SimTradeCollector()
	// Subscribe all target
	choosetarget.SubscribeTarget(&global.TargetArr)
	// list targets
	for i, v := range global.TargetArr {
		logger.GetLogger().WithFields(map[string]interface{}{
			"Date": global.LastTradeDay.Format(global.ShortTimeLayout),
			"Rank": i + 1,
			"Name": global.AllStockNameMap.GetValueByKey(v),
		}).Infof("Volume Rank")
	}
	// fetch entiretick
	logger.GetLogger().Info("FetchEntireTick and FetchKbar")
	fetchentiretick.FetchEntireTick(global.TargetArr, []time.Time{global.LastTradeDay}, global.BaseCond)
	// Fetch Kbar
	kbarTradeDayArr, err := importbasic.GetLastNTradeDay(sysparminit.GlobalSettings.GetKbarPeriod())
	if err != nil {
		logger.GetLogger().Panic(err)
	}
	fetchentiretick.GetAndSaveKbar(global.TargetArr, kbarTradeDayArr[len(kbarTradeDayArr)-1], kbarTradeDayArr[0])
	logger.GetLogger().Info("FetchEntireTick and FetchKbar Done")
	// Background get trade record
	go tradebot.CheckOrderStatus()
	// Init quota and tradeday order map
	go tradebot.InitStartUpQuota()
	// Monitor TSE001 Status
	go choosetarget.TSEProcess()
	// Add Top Rank Targets
	go choosetarget.AddTop10RankTarget()
	logger.GetLogger().Info("TradeProcess Success Started")
}

func simulatationEntry() {
	var err error
	tmpChan := make(chan string)
	deployment := os.Getenv("DEPLOYMENT")
	if deployment != "docker" {
		prompt := promptui.Prompt{
			Label: "Simulate?(y/n)",
		}
		var result string
		result, err = prompt.Run()
		if err != nil {
			logger.GetLogger().Panic(err)
		}
		if result == "y" {
			ansArr := simulationPrompt()
			simulateprocess.Simulate(ansArr[0], ansArr[1], ansArr[2], ansArr[3])
			<-tmpChan
		}
	}
	getConds()
}

func getConds() {
	var err error
	tmpChan := make(chan string)
	global.ForwardCond, err = simulate.GetBestForwardCondByTradeDay(global.TradeDay, database.GetAgent())
	if err != nil {
		logger.GetLogger().Panic(err)
	}
	if global.ForwardCond.Model.ID == 0 {
		simulateprocess.Simulate("a", "n", "n", "1")
		global.ForwardCond, err = simulate.GetBestForwardCondByTradeDay(global.TradeDay, database.GetAgent())
		if err != nil {
			logger.GetLogger().Panic(err)
		}
	}
	global.ReverseCond, err = simulate.GetBestReverseCondByTradeDay(global.TradeDay, database.GetAgent())
	if err != nil {
		logger.GetLogger().Panic(err)
	}
	if global.ReverseCond.Model.ID == 0 {
		simulateprocess.Simulate("b", "n", "n", "1")
		global.ReverseCond, err = simulate.GetBestReverseCondByTradeDay(global.TradeDay, database.GetAgent())
		if err != nil {
			logger.GetLogger().Panic(err)
		}
	}
	if global.ForwardCond.Model.ID == 0 || global.ReverseCond.Model.ID == 0 {
		logger.GetLogger().Warn("no cond to trade")
		<-tmpChan
	}
	// forwardResult, err := simulate.GetResultByCond(int(global.ForwardCond.ID), database.GetAgent())
	// if err != nil {
	// 	logger.GetLogger().Panic(err)
	// }
	// reverseResult, err := simulate.GetResultByCond(int(global.ReverseCond.ID), database.GetAgent())
	// if err != nil {
	// 	logger.GetLogger().Panic(err)
	// }
	// if float64(forwardResult.Balance)/float64(reverseResult.Balance) < 0.5 {
	// 	global.TradeSwitch.Buy = false
	// 	logger.GetLogger().Warn("TradeSwitch Buy is OFF")
	// }
	// if float64(forwardResult.Balance)/float64(reverseResult.Balance) > 2 {
	// 	global.TradeSwitch.SellFirst = false
	// 	logger.GetLogger().Warn("TradeSwitch SellFirst is OFF")
	// }
	simulateprocess.ClearAllNotBest()
}

func simulationPrompt() []string {
	prompt := promptui.Prompt{
		Label: "Balance type?(a: forward, b: reverse)",
	}
	balanceTypeAns, err := prompt.Run()
	if err != nil {
		logger.GetLogger().Panic(err)
	}
	prompt = promptui.Prompt{
		Label: "Discard over time trade?(y/n)",
	}
	discardAns, err := prompt.Run()
	if err != nil {
		logger.GetLogger().Panic(err)
	}
	prompt = promptui.Prompt{
		Label: "Use global cond?(y/n)",
	}
	useDefault, err := prompt.Run()
	if err != nil {
		logger.GetLogger().Panic(err)
	}
	prompt = promptui.Prompt{
		Label: "N days?",
	}
	countAns, err := prompt.Run()
	if err != nil {
		logger.GetLogger().Panic(err)
	}
	return []string{balanceTypeAns, discardAns, useDefault, countAns}
}
