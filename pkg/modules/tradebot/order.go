// Package tradebot package tradebot
package tradebot

import (
	"errors"
	"net/http"

	"gitlab.tocraw.com/root/toc_trader/external/sinopacsrv"
	"gitlab.tocraw.com/root/toc_trader/internal/restful"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
)

// PlaceOrder PlaceOrder
func PlaceOrder(action OrderAction, stockNum string, stockQuantity int64, stockPrice float64) (returnOrder sinopacsrv.OrderResponse, err error) {
	var url string
	switch action {
	case BuyAction:
		url = "/pyapi/trade/buy"
	case SellAction:
		url = "/pyapi/trade/sell"
	case SellFirstAction:
		url = "/pyapi/trade/sell_first"
	}
	order := OrderBody{
		Stock:    stockNum,
		Price:    stockPrice,
		Quantity: stockQuantity,
	}
	resp, err := restful.GetClient().R().
		SetBody(order).
		SetResult(&sinopacsrv.OrderResponse{}).
		Post("http://" + global.PyServerHost + ":" + global.PyServerPort + url)
	if err != nil {
		return returnOrder, err
	} else if resp.StatusCode() != http.StatusOK {
		return returnOrder, errors.New("PlaceOrder api fail")
	}
	res := *resp.Result().(*sinopacsrv.OrderResponse)
	return res, err
}

// Cancel Cancel
func Cancel(orderID string) (err error) {
	order := CancelBody{
		OrderID: orderID,
	}
	resp, err := restful.GetClient().R().
		SetBody(order).
		SetResult(&sinopacsrv.OrderResponse{}).
		Post("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/trade/cancel")
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("Cancel api fail")
	}
	res := *resp.Result().(*sinopacsrv.OrderResponse)
	switch res.Status {
	case sinopacsrv.StatusFail:
		return errors.New(sinopacsrv.StatusFail)
	case sinopacsrv.StatusAlreadyCanceled:
		return errors.New(sinopacsrv.StatusAlreadyCanceled)
	case sinopacsrv.StatusCancelOrderNotFound:
		return errors.New(sinopacsrv.StatusCancelOrderNotFound)
	}
	return err
}

// FetchOrderStatus FetchOrderStatus
func FetchOrderStatus() (err error) {
	resp, err := restful.GetClient().R().
		SetResult(&sinopacsrv.SinoStatusResponse{}).
		Get("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/trade/status")
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("FetchOrderStatus api fail")
	}
	res := *resp.Result().(*sinopacsrv.SinoStatusResponse)
	if res.Status != sinopacsrv.StatusSuccuss {
		return errors.New("FetchOrderStatus fail")
	}
	return err
}
