// Package tradebot package tradebot
package tradebot

import (
	"errors"

	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/pyresponse"
)

// PlaceOrder PlaceOrder
func PlaceOrder(action OrderAction, stockNum string, stockQuantity int64, stockPrice float64) (returnOrder pyresponse.PyServerResponse, err error) {
	if stockNum == "" || stockQuantity == 0 {
		return returnOrder, errors.New("PlaceOrder input error")
	}
	var url string
	if action == BuyAction {
		url = "/pyapi/trade/buy"
	} else if action == SellAction {
		url = "/pyapi/trade/sell"
	}
	order := OrderBody{
		Stock:    stockNum,
		Price:    stockPrice,
		Quantity: stockQuantity,
	}
	resp, err := global.RestyClient.R().
		SetBody(order).
		SetResult(&pyresponse.PyServerResponse{}).
		Post("http://" + global.PyServerHost + ":" + global.PyServerPort + url)
	if err != nil {
		return returnOrder, err
	} else if resp.StatusCode() != 200 {
		return returnOrder, errors.New("api fail")
	}
	res := *resp.Result().(*pyresponse.PyServerResponse)
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
	resp, err := global.RestyClient.R().
		SetBody(order).
		SetResult(&pyresponse.PyServerResponse{}).
		Post("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/trade/cancel")
	if err != nil {
		return err
	} else if resp.StatusCode() != 200 {
		return errors.New("api fail")
	}
	res := *resp.Result().(*pyresponse.PyServerResponse)
	if res.Status == "fail" {
		return errors.New("cancel fail")
	}
	return err
}
