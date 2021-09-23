// Package core package core
package core

import (
	"errors"
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
	// Development
	deployment := os.Getenv("DEPLOYMENT")
	if deployment != "docker" {
		// Send ip to sinopac srv
		sendCurrentIP()
		// Simulate
		simulate.Simulate()
	}
	// Generate global target array
	choosetarget.GetTarget(sysparminit.GlobalSettings.GetTargetCondArr())
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

		tick := time.NewTicker(10 * time.Second)
		for range tick.C {
			var count int
			if newTargetArr, err := choosetarget.GetTopTarget(15); err != nil {
				logger.Logger.Error(err)
				continue
			} else {
				count = len(newTargetArr)
				if len(newTargetArr) != 0 {
					choosetarget.SubscribeTarget(newTargetArr)
					global.TargetArr = append(global.TargetArr, newTargetArr...)
				}
			}
			if count != 0 {
				logger.Logger.Infof("GetTopTarget %d", count)
			}
		}
	}
}

func sendCurrentIP() {
	var err error
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
