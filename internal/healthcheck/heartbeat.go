// Package healthcheck package healthcheck
package healthcheck

import (
	"errors"
	"net/http"

	"gitlab.tocraw.com/root/toc_trader/external/sinopacsrv"
	"gitlab.tocraw.com/root/toc_trader/internal/rest"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
)

// FullRestart FullRestart
func FullRestart() (err error) {
	if err = askSinopacSRVRestart(); err != nil {
		return err
	}
	RestartService()
	return err
}

func askSinopacSRVRestart() error {
	resp, err := rest.GetClient().R().
		SetResult(&sinopacsrv.OrderResponse{}).
		Get("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/system/restart")
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("askSinopacSRVRestart api fail")
	}
	res := *resp.Result().(*sinopacsrv.OrderResponse)
	if res.Status != "success" {
		return errors.New(res.Status)
	}
	return err
}

// RestartService RestartService
func RestartService() {
	global.ExitChannel <- global.ExitSignal
}