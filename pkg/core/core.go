// Package core package core
package core

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/manifoldco/promptui"
	"gitlab.tocraw.com/root/toc_trader/init/sysparminit"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/pyresponse"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/stock"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/targetstock"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/choosetarget"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/fetchentiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/importbasic"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/simulateprocess"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tradebot"
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
			fmt.Printf("%s volume rank no. %d is %s\n", global.LastTradeDay.Format(global.ShortTimeLayout), i+1, global.AllStockNameMap.GetName(v))
		}
		tmp := []time.Time{global.LastTradeDay}
		fetchentiretick.FetchEntireTick(targets, tmp, global.TickAnalyzeCondition)
		if dbTarget, err := targetstock.GetTargetByTime(global.LastTradeDay, global.GlobalDB); err != nil {
			panic(err)
		} else if len(dbTarget) == 0 {
			logger.Logger.Info("Saving targets")
			targetStockArr, err := stock.GetStocksFromNumArr(targets, global.GlobalDB)
			if err != nil {
				panic(err)
			}
			for _, v := range targetStockArr {
				savedTarget = append(savedTarget, targetstock.Target{
					LastTradeDay: global.LastTradeDay,
					Stock:        v,
				})
			}
			if err := targetstock.InsertMultiTarget(savedTarget, global.GlobalDB); err != nil {
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
	logger.Logger.Info("FetchEntireTick and Kbar Done")

	// Check tradeday or target exist
	if len(global.LastTradeDayArr) == 0 || len(global.TargetArr) == 0 {
		logger.Logger.Warn("no trade day or no target")
	} else {
		// Background get trade record
		logger.Logger.Info("Background tasks starts")
		go tradebot.CheckOrderStatusLoop()
		// Monitor TSE001 Status
		go choosetarget.TSEProcess()
		// Add Top Rank Targets
		go addRankTarget()
	}
}

func sendCurrentIP() {
	var err error
	results := findMachineIP()
	resp, err := global.RestyClient.R().
		SetHeader("X-Trade-Bot-Host", results[len(results)-1]).
		SetResult(&pyresponse.PyServerResponse{}).
		Post("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/system/tradebothost")
	if err != nil {
		logger.Logger.Error(err)
		return
	} else if resp.StatusCode() != 200 {
		logger.Logger.Error("SendCurrentIP api fail")
		return
	}
	res := *resp.Result().(*pyresponse.PyServerResponse)
	if res.Status != "success" {
		logger.Logger.Error(errors.New("sendCurrentIP fail"))
	}
}

func findMachineIP() []string {
	var results []string
	ifaces, err := net.Interfaces()
	if err != nil {
		logger.Logger.Error(err)
	}
	var addrs []net.Addr
	for _, i := range ifaces {
		if i.HardwareAddr.String() == "" {
			continue
		}
		addrs, err = i.Addrs()
		if err != nil {
			logger.Logger.Error(err)
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

func checkIsOpenTime() bool {
	starTime := global.TradeDay.Add(1 * time.Hour)
	if time.Now().After(starTime) && time.Now().Before(global.TradeDayEndTime) {
		return true
	}
	return false
}
