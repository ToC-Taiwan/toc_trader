// Package core package core
package core

import (
	"errors"
	"net"

	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
	"gitlab.tocraw.com/root/toc_trader/tools/rest"
)

func sendCurrentIP() {
	var err error
	results := findMachineIP()
	resp, err := rest.GetClient().R().
		SetHeader("X-Trade-Bot-Host", results[len(results)-1]).
		SetResult(&global.PyServerResponse{}).
		Post("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/system/tradebothost")
	if err != nil {
		logger.GetLogger().Error(err)
		return
	} else if resp.StatusCode() != 200 {
		logger.GetLogger().Error("SendCurrentIP api fail")
		return
	}
	res := *resp.Result().(*global.PyServerResponse)
	if res.Status != "success" {
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
