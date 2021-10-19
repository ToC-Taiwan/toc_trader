// Package healthcheck package healthcheck
package healthcheck

import (
	"errors"

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
	resp, err := global.RestyClient.R().
		SetResult(&global.PyServerResponse{}).
		Get("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/system/restart")
	if err != nil {
		return err
	} else if resp.StatusCode() != 200 {
		return errors.New("askSinopacSRVRestart api fail")
	}
	res := *resp.Result().(*global.PyServerResponse)
	if res.Status != "success" {
		return errors.New(res.Status)
	}
	return err
}

// RestartService RestartService
func RestartService() {
	global.ExitChannel <- global.ExitSignal
}
