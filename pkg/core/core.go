// Package core package core
package core

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"gitlab.tocraw.com/root/toc_trader/init/sysparminit"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/pyresponse"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/choosetarget"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/fetchentiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/importbasic"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/simulate"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tradebot"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
)

// TradeProcess TradeProcess
func TradeProcess() {
	// Import all stock and update AllStockNameMap
	importbasic.ImportAllStock()
	// If needed, it will update StockCloseMap and volume, close in database
	choosetarget.UpdateLastStockVolume()
	// Monitor TSE001 Status
	go choosetarget.TSEProcess()
	// Development
	deployment := os.Getenv("DEPLOYMENT")
	if deployment != "docker" {
		// Send ip to sinopac srv
		sendCurrentIP()
		// Simulate
		fmt.Print("Need simulate?(y/n): ")
		reader := bufio.NewReader(os.Stdin)
		ans, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		if ans == "y\n" {
			logger.Logger.Warn("Simulating")
			simulate.Simulate()
		}
	}
	// Generate global target array
	choosetarget.GetTargetFromStockList(sysparminit.GlobalSettings.GetTargetCondArr())
	// UnSubscribeAll first
	choosetarget.UnSubscribeAll()

	// Check tradeday or target exist
	if len(global.LastTradeDayArr) == 0 || len(global.TargetArr) == 0 {
		logger.Logger.Warn("no trade day or no target")
	} else {
		// Subscribe all target
		choosetarget.SubscribeTarget(global.TargetArr)
		// Put data into channel, it will wait for all fetch done. Final close all channel.
		fetchentiretick.FetchEntireTick(global.TargetArr, global.LastTradeDayArr, global.TickAnalyzeCondition)
		logger.Logger.Info("FetchEntireTick Done")
		// Background get trade record
		go tradebot.CheckOrderStatusLoop()
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
