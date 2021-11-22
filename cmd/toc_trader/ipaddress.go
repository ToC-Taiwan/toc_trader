// Package main package main
package main

import (
	"net"
	"net/http"

	"gitlab.tocraw.com/root/toc_trader/external/sinopacsrv"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/internal/restful"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
)

func sendHostIP(ip string) error {
	var err error
	resp, err := restful.GetClient().R().
		SetHeader("X-Trade-Bot-Host", ip).
		SetResult(&sinopacsrv.OrderResponse{}).
		Post("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/system/tradebothost")
	if err != nil {
		logger.GetLogger().Error(err)
		return err
	} else if resp.StatusCode() != http.StatusOK {
		logger.GetLogger().Error("SendCurrentIP api fail")
		return err
	}
	res := *resp.Result().(*sinopacsrv.OrderResponse)
	if res.Status != sinopacsrv.StatusSuccuss {
		logger.GetLogger().Error(res.Status)
	}
	return err
}

func getHostIP() string {
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
	return results[len(results)-1]
}
