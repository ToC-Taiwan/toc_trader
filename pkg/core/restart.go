// Package core package core
package core

import (
	"errors"

	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/pyresponse"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/process"
)

// FullRestart FullRestart
func FullRestart() (err error) {
	if err = askSinopacSRVRestart(); err != nil {
		return err
	}
	process.RestartService()
	return err
}

func askSinopacSRVRestart() error {
	resp, err := global.RestyClient.R().
		SetResult(&pyresponse.PyServerResponse{}).
		Get("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/system/restart")
	if err != nil {
		return err
	} else if resp.StatusCode() != 200 {
		return errors.New("askSinopacSRVRestart api fail")
	}
	res := *resp.Result().(*pyresponse.PyServerResponse)
	if res.Status != "success" {
		return errors.New("askSinopacSRVRestart fail")
	}
	return err
}
