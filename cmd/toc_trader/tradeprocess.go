// Package main package main
package main

import (
	"errors"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/manifoldco/promptui"
	"gitlab.tocraw.com/root/toc_trader/external/sinopacsrv"
	"gitlab.tocraw.com/root/toc_trader/init/sysparminit"
	"gitlab.tocraw.com/root/toc_trader/internal/database"
	"gitlab.tocraw.com/root/toc_trader/internal/healthcheck"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/internal/restful"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
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
	// Check token is expired or not, if expired, restart service
	go checkSinopacToken()
	// Import all stock and update AllStockNameMap
	importbasic.ImportAllStock()
	// Chekc if is in debug mode and simulation
	simulatationEntry()

	// Generate global target array and fetch entireTick
	var savedTarget []targetstock.Target
	if targets, err := choosetarget.GetTargetByVolumeRankByDate(global.LastTradeDay.Format(global.ShortTimeLayout), 200); err != nil {
		panic(err)
	} else {
		// UnSubscribeAll first
		choosetarget.UnSubscribeAll()
		// Clear stram tick in database
		if err := tickprocess.DeleteAllStreamTicks(); err != nil {
			panic(err)
		}
		// Simtrade collect
		go subscribe.SimTradeCollector(len(targets) * 2)
		// Subscribe all target
		global.TargetArr = targets
		choosetarget.SubscribeTarget(global.TargetArr)
		for i, v := range targets {
			logger.GetLogger().WithFields(map[string]interface{}{
				"Date": global.LastTradeDay.Format(global.ShortTimeLayout),
				"Rank": i + 1,
				"Name": global.AllStockNameMap.GetName(v),
			}).Infof("Volume Rank")
		}
		tmp := []time.Time{global.LastTradeDay}
		fetchentiretick.FetchEntireTick(targets, tmp, global.BaseCond)
		if dbTarget, err := targetstock.GetTargetByTime(global.LastTradeDay, database.GetAgent()); err != nil {
			panic(err)
		} else if len(dbTarget) == 0 {
			logger.GetLogger().Info("Saving targets")
			targetStockArr, err := stock.GetStocksFromNumArr(targets, database.GetAgent())
			if err != nil {
				panic(err)
			}
			for _, v := range targetStockArr {
				savedTarget = append(savedTarget, targetstock.Target{
					LastTradeDay: global.LastTradeDay,
					Stock:        v,
				})
			}
			if err := targetstock.InsertMultiTarget(savedTarget, database.GetAgent()); err != nil {
				panic(err)
			}
		}
	}
	// Fetch Kbar
	kbarTradeDayArr, err := importbasic.GetLastNTradeDay(sysparminit.GlobalSettings.GetKbarPeriod())
	if err != nil {
		logger.GetLogger().Error(err)
	}
	fetchentiretick.FetchKbar(global.TargetArr, kbarTradeDayArr[len(kbarTradeDayArr)-1], kbarTradeDayArr[0])
	logger.GetLogger().Info("FetchEntireTick and Kbar Done")

	// Check tradeday or target exist
	if len(global.LastTradeDayArr) == 0 || len(global.TargetArr) == 0 {
		panic("no trade day or no target")
	} else {
		// Background get trade record
		logger.GetLogger().Info("Background Tasks Start")
		go tradebot.CheckOrderStatusLoop()
		// Init quota and tradeday order map
		go tradebot.InitStartUpQuota()
		// Monitor TSE001 Status
		go choosetarget.TSEProcess()
		// Add Top Rank Targets
		go addRankTarget()
	}
}

func checkSinopacToken() {
	for range time.Tick(10 * time.Second) {
		if err := healthcheck.CheckSinopacSRVStatus(); err != nil {
			panic(err)
		}
	}
}

func simulatationEntry() {
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
			simulateprocess.ClearAllSimulation()
			prompt := promptui.Prompt{
				Label: "Balance type?(a: forward, b: reverse)",
			}
			balanceTypeAns, err := prompt.Run()
			if err != nil {
				panic(err)
			}
			prompt = promptui.Prompt{
				Label: "Discard over time trade?(y/n)",
			}
			discardAns, err := prompt.Run()
			if err != nil {
				panic(err)
			}
			prompt = promptui.Prompt{
				Label: "Use global cond?(y/n)",
			}
			useGlobalAns, err := prompt.Run()
			if err != nil {
				panic(err)
			}
			prompt = promptui.Prompt{
				Label: "N days?",
			}
			countAns, err := prompt.Run()
			if err != nil {
				panic(err)
			}
			simulateprocess.Simulate(balanceTypeAns, discardAns, useGlobalAns, countAns)
		} else {
			tmpChan := make(chan string)
			logger.GetLogger().Warn("Please run in container mode")
			tmpChan <- "stockHere"
		}
	} else {
		var err error
		global.ForwardCond, err = simulate.GetBestForwardCond(database.GetAgent())
		if err != nil {
			panic(err)
		}
		if global.ForwardCond.Model.ID == 0 {
			simulateprocess.Simulate("a", "n", "n", "1")
			global.ForwardCond, err = simulate.GetBestForwardCond(database.GetAgent())
			if err != nil {
				panic(err)
			}
		}
		global.ReverseCond, err = simulate.GetBestReverseCond(database.GetAgent())
		if err != nil {
			panic(err)
		}
		if global.ReverseCond.Model.ID == 0 {
			simulateprocess.Simulate("b", "n", "n", "1")
			global.ReverseCond, err = simulate.GetBestReverseCond(database.GetAgent())
			if err != nil {
				panic(err)
			}
		}
		if global.ForwardCond.Model.ID == 0 || global.ReverseCond.Model.ID == 0 {
			panic("no cond to trade")
		}
		logger.GetLogger().Warnf("BestForward is %+v", global.ForwardCond)
		logger.GetLogger().Warnf("BestReverse is %+v", global.ReverseCond)
	}
}

func sendCurrentIP() {
	var err error
	results := findMachineIP()
	resp, err := restful.GetClient().R().
		SetHeader("X-Trade-Bot-Host", results[len(results)-1]).
		SetResult(&sinopacsrv.OrderResponse{}).
		Post("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/system/tradebothost")
	if err != nil {
		logger.GetLogger().Error(err)
		return
	} else if resp.StatusCode() != http.StatusOK {
		logger.GetLogger().Error("SendCurrentIP api fail")
		return
	}
	res := *resp.Result().(*sinopacsrv.OrderResponse)
	if res.Status != sinopacsrv.StatusSuccuss {
		logger.GetLogger().Error(errors.New("sendCurrentIP fail"))
	}
}

func findMachineIP() []string {
	var results []string
	ifaces, err := net.Interfaces()
	if err != nil {
		logger.GetLogger().Error(err)
	}
	var addrs []net.Addr
	for _, i := range ifaces {
		if i.HardwareAddr.String() == "" {
			continue
		}
		addrs, err = i.Addrs()
		if err != nil {
			logger.GetLogger().Error(err)
		}
		for _, addr := range addrs {
			if ip := addr.(*net.IPNet).IP.To4(); ip != nil {
				if ip[0] != 127 && ip[0] != 169 {
					results = append(results, ip.String())
				}
			}
		}
	}
	return results
}

func addRankTarget() {
	tick := time.Tick(15 * time.Second)
	for range tick {
		if !tradebot.CheckIsOpenTime() {
			continue
		}
		var count int
		if newTargetArr, err := choosetarget.GetTopTarget(20); err != nil {
			logger.GetLogger().Error(err)
			continue
			// Start from 9:10
		} else if time.Now().After(global.TradeDay.Add(1*time.Hour + 10*time.Minute)) {
			count = len(newTargetArr)
			if count != 0 {
				choosetarget.SubscribeTarget(newTargetArr)
				global.TargetArr = append(global.TargetArr, newTargetArr...)
			}
		}
		if count != 0 {
			logger.GetLogger().Infof("GetTopTarget %d", count)
		}
	}
}
