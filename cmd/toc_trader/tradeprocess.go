// Package main package main
package main

import (
	"time"

	"gitlab.tocraw.com/root/toc_trader/init/sysparminit"
	"gitlab.tocraw.com/root/toc_trader/internal/database"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/stock"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/targetstock"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/choosetarget"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/fetchentiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/importbasic"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/subscribe"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tickprocess"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tradebot"
)

// TradeProcess TradeProcess
func TradeProcess() {
	// Check tradeday or target exist
	if len(global.LastTradeDayArr) == 0 {
		logger.GetLogger().Panic("No trade day")
	}
	// Import all stock and update AllStockNameMap
	importbasic.ImportAllStock()
	// Clear stram tick in database
	if err := tickprocess.DeleteAllStreamTicks(); err != nil {
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
	choosetarget.UnSubscribeAll()
	// Simtrade collect
	go subscribe.SimTradeCollector()
	// Subscribe all target
	choosetarget.SubscribeTarget(&global.TargetArr)
	// list targets
	for i, v := range global.TargetArr {
		logger.GetLogger().WithFields(map[string]interface{}{
			"Date": global.LastTradeDay.Format(global.ShortTimeLayout),
			"Rank": i + 1,
			"Name": global.AllStockNameMap.GetName(v),
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
	fetchentiretick.FetchKbar(global.TargetArr, kbarTradeDayArr[len(kbarTradeDayArr)-1], kbarTradeDayArr[0])
	logger.GetLogger().Info("FetchEntireTick and FetchKbar Done")
	// Background get trade record
	go tradebot.CheckOrderStatusLoop()
	// Init quota and tradeday order map
	go tradebot.InitStartUpQuota()
	// Monitor TSE001 Status
	go choosetarget.TSEProcess()
	// Add Top Rank Targets
	go choosetarget.AddTop10RankTarget()
	logger.GetLogger().Info("TradeProcess Success Started")
}
