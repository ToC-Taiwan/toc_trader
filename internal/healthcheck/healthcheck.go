// Package healthcheck package healthcheck
package healthcheck

import (
	"errors"
	"net/http"

	"gitlab.tocraw.com/root/toc_trader/external/sinopacsrv"
	"gitlab.tocraw.com/root/toc_trader/internal/restful"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
)

// AskSinopacSRVRestart AskSinopacSRVRestart
func AskSinopacSRVRestart() error {
	resp, err := restful.GetClient().R().
		SetResult(&sinopacsrv.OrderResponse{}).
		Get("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/system/restart")
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("AskSinopacSRVRestart api fail")
	}
	res := *resp.Result().(*sinopacsrv.OrderResponse)
	if res.Status != "success" {
		return errors.New(res.Status)
	}
	return err
}

// ExitService ExitService
func ExitService() {
	global.ExitChannel <- global.ExitSignal
}

// GetSinopacSRVToken GetSinopacSRVToken
func GetSinopacSRVToken() (token string, err error) {
	resp, err := restful.GetClient().R().
		SetResult(&sinopacsrv.SinopacHealthStatus{}).
		Get("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/system/healthcheck")
	if err != nil {
		return token, err
	} else if resp.StatusCode() != http.StatusOK {
		return token, errors.New("GetSinopacSRVToken api fail")
	}
	res := *resp.Result().(*sinopacsrv.SinopacHealthStatus)
	return res.ServerToken, err
}
