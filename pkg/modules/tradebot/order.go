// Package tradebot package tradebot
package tradebot

import (
	"errors"

	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/tools/rest"
)

// PlaceOrder PlaceOrder
func PlaceOrder(action OrderAction, stockNum string, stockQuantity int64, stockPrice float64) (returnOrder global.PyServerResponse, err error) {
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
	resp, err := rest.GetClient().R().
		SetBody(order).
		SetResult(&global.PyServerResponse{}).
		Post("http://" + global.PyServerHost + ":" + global.PyServerPort + url)
	if err != nil {
		return returnOrder, err
	} else if resp.StatusCode() != 200 {
		return returnOrder, errors.New("PlaceOrder api fail")
	}
	res := *resp.Result().(*global.PyServerResponse)
	return res, err
}

// Cancel Cancel
func Cancel(orderID string) (err error) {
	if orderID == "" {
		return errors.New("Cancel input error")
	}
	order := CancelBody{
		OrderID: orderID,
	}
	resp, err := rest.GetClient().R().
		SetBody(order).
		SetResult(&global.PyServerResponse{}).
		Post("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/trade/cancel")
	if err != nil {
		return err
	} else if resp.StatusCode() != 200 {
		return errors.New("Cancel api fail")
	}
	res := *resp.Result().(*global.PyServerResponse)
	if res.Status == "fail" {
		return errors.New("Cancel fail")
	}
	return err
}
