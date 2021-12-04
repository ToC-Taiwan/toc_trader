// Package sinopacapi package sinopacapi
package sinopacapi

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"google.golang.org/protobuf/proto"
)

var client *TradeAgent

// TradeAgent TradeAgent
type TradeAgent struct {
	restAgent *resty.Client
	urlPrefix string
}

// Order Order
type Order struct {
	StockNum  string
	Price     float64
	Quantity  int64
	OrderID   string
	Action    OrderAction
	TradeTime time.Time
}

// GetAgent GetAgent
func GetAgent() *TradeAgent {
	if client == nil {
		panic("trade agent not initital")
	}
	return client
}

// NewAgent NewAgent
func NewAgent(serverHost, serverPort string, restAgent *resty.Client) {
	new := TradeAgent{
		restAgent: restAgent,
		urlPrefix: "http://" + serverHost + ":" + serverPort,
	}
	client = &new
}

// FetchServerKey FetchServerKey
func (c *TradeAgent) FetchServerKey() (token string, err error) {
	var resp *resty.Response
	resp, err = c.restAgent.R().
		SetResult(&ResponseHealthStatus{}).
		Get(c.urlPrefix + urlFetchServerKey)
	if err != nil {
		return token, err
	} else if resp.StatusCode() != http.StatusOK {
		return token, errors.New("FetchServerKey API Fail")
	}
	return resp.Result().(*ResponseHealthStatus).ServerToken, err
}

// UpdateTraderIP UpdateTraderIP
func (c *TradeAgent) UpdateTraderIP(ip string) (err error) {
	var resp *resty.Response
	resp, err = c.restAgent.R().
		SetHeader("X-Trade-Bot-Host", ip).
		SetResult(&ResponseCommon{}).
		Post(c.urlPrefix + urlUpdateTraderIP)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("sendHostIP API Fail")
	}
	if status := resp.Result().(*ResponseCommon).Status; status != StatusSuccuss {
		return errors.New(status)
	}
	return err
}

// PlaceOrder PlaceOrder
func (c *TradeAgent) PlaceOrder(order Order) (res OrderResponse, err error) {
	var url string
	switch order.Action {
	case ActionBuy:
		url = urlPlaceOrderBuy
	case ActionSell:
		url = urlPlaceOrderSell
	case ActionSellFirst:
		url = urlPlaceOrderSellFirst
	}
	body := OrderBody{
		Stock:    order.StockNum,
		Price:    order.Price,
		Quantity: order.Quantity,
	}
	var resp *resty.Response
	resp, err = c.restAgent.R().
		SetBody(body).
		SetResult(&OrderResponse{}).
		Post(c.urlPrefix + url)
	if err != nil {
		return res, err
	} else if resp.StatusCode() != http.StatusOK {
		return res, errors.New("PlaceOrder API Fail")
	}
	return *resp.Result().(*OrderResponse), err
}

// CancelOrder CancelOrder
func (c *TradeAgent) CancelOrder(orderID string) (err error) {
	order := OrderCancelBody{
		OrderID: orderID,
	}
	var resp *resty.Response
	resp, err = c.restAgent.R().
		SetBody(order).
		SetResult(&OrderResponse{}).
		Post(c.urlPrefix + urlCancelOrder)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("CancelOrder API Fail")
	}
	switch resp.Result().(*OrderResponse).Status {
	case StatusFail:
		return errors.New(StatusFail)
	case StatusAlreadyCanceled:
		return errors.New(StatusAlreadyCanceled)
	case StatusCancelOrderNotFound:
		return errors.New(StatusCancelOrderNotFound)
	}
	return err
}

// FetchOrderStatus FetchOrderStatus
func (c *TradeAgent) FetchOrderStatus() (err error) {
	var resp *resty.Response
	resp, err = c.restAgent.R().
		SetResult(&ResponseCommon{}).
		Get(c.urlPrefix + urlFetchOrderStatus)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("FetchOrderStatus API Fail")
	}
	if status := resp.Result().(*ResponseCommon).Status; status != StatusSuccuss {
		return errors.New(status)
	}
	return err
}

// RestartSinopacSRV RestartSinopacSRV
func (c *TradeAgent) RestartSinopacSRV() (err error) {
	var resp *resty.Response
	resp, err = c.restAgent.R().
		SetResult(&OrderResponse{}).
		Get(c.urlPrefix + urlRestartSinopacSRV)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("RestartSinopacSRV API Fail")
	}
	if status := resp.Result().(*OrderResponse).Status; status != StatusSuccuss {
		return errors.New(status)
	}
	return err
}

// FetchLastCloseByStockArrDateArr FetchLastCloseByStockArrDateArr
func (c *TradeAgent) FetchLastCloseByStockArrDateArr(stockNumArr, dateArr []string) (stockLastCloseMap map[string][]float64, err error) {
	stockLastCloseMap = make(map[string][]float64)
	stockAndDateArr := FetchLastCloseBody{
		StockNumArr: stockNumArr,
		DateArr:     dateArr,
	}
	var resp *resty.Response
	resp, err = c.restAgent.R().
		SetBody(stockAndDateArr).
		SetResult(&[]LastCloseWithStockAndDate{}).
		Post(c.urlPrefix + urlFetchLastCloseByStockArrDateArr)
	if err != nil {
		return stockLastCloseMap, err
	} else if resp.StatusCode() != http.StatusOK {
		return stockLastCloseMap, errors.New("FetchLastCloseByStockArrDateArr API Fail")
	}
	res := *resp.Result().(*[]LastCloseWithStockAndDate)
	for _, v := range res {
		var tmp []float64
		for k := len(v.CloseArr) - 1; k >= 0; k-- {
			tmp = append(tmp, v.CloseArr[k].Close)
		}
		stockLastCloseMap[v.StockNum] = tmp
	}
	return stockLastCloseMap, err
}

// FetchAllSnapShot FetchAllSnapShot
func (c *TradeAgent) FetchAllSnapShot() (data []*SnapShotProto, err error) {
	var resp *resty.Response
	resp, err = c.restAgent.R().
		Get(c.urlPrefix + urlFetchAllSnapShot)
	if err != nil {
		return data, err
	} else if resp.StatusCode() != http.StatusOK {
		return data, errors.New("FetchAllSnapShot API Fail")
	}
	body := SnapShotArrProto{}
	if err = proto.Unmarshal(resp.Body(), &body); err != nil {
		return data, err
	}
	return body.Data, err
}

// FetchStockCloseMapByStockDateArr FetchStockCloseMapByStockDateArr will return map[stock][date]close
func (c *TradeAgent) FetchStockCloseMapByStockDateArr(stockNumArr []string, dateArr []time.Time) (result map[string]map[string]float64, noCloseArr []string, err error) {
	result = make(map[string]map[string]float64)
	stockArr := FetchLastCountBody{
		StockNumArr: stockNumArr,
	}
	tmpNoClose := make(map[string]bool)
	for _, date := range dateArr {
		var resp *resty.Response
		resp, err = c.restAgent.R().
			SetHeader("X-Date", date.Format(shortTimeLayout)).
			SetBody(stockArr).
			SetResult(&[]StockLastCount{}).
			Post(c.urlPrefix + urlFetchStockCloseMapByStockDateArr)
		if err != nil {
			return result, noCloseArr, err
		} else if resp.StatusCode() != http.StatusOK {
			return result, noCloseArr, errors.New("FetchStockCloseMapByStockDateArr API Fail")
		}
		stockLastCountArr := *resp.Result().(*[]StockLastCount)
		for _, val := range stockLastCountArr {
			if len(val.Close) != 0 {
				if result[val.Code] == nil {
					result[val.Code] = make(map[string]float64)
				}
				result[val.Code][val.Date] = val.Close[0]
			} else {
				delete(result, val.Code)
				if _, ok := tmpNoClose[val.Code]; !ok {
					tmpNoClose[val.Code] = true
				}
			}
		}
	}
	for k := range tmpNoClose {
		noCloseArr = append(noCloseArr, k)
	}
	return result, noCloseArr, err
}

// FetchTSE001CloseByDate FetchTSE001CloseByDate
func (c *TradeAgent) FetchTSE001CloseByDate(date time.Time) (close float64, err error) {
	var resp *resty.Response
	resp, err = c.restAgent.R().
		SetHeader("X-Date", date.Format(shortTimeLayout)).
		SetResult(&[]StockLastCount{}).
		Post(c.urlPrefix + urlFetchTSE001CloseByDate)
	if err != nil {
		return close, err
	} else if resp.StatusCode() != http.StatusOK {
		return close, errors.New("FetchTSE001CloseByDate API Fail")
	}
	stockLastCountArr := *resp.Result().(*[]StockLastCount)
	return stockLastCountArr[0].Close[0], err
}

// FetchVolumeRankByDate FetchVolumeRankByDate
func (c *TradeAgent) FetchVolumeRankByDate(date string, count int64) (data []*VolumeRankProto, err error) {
	var resp *resty.Response
	resp, err = c.restAgent.R().
		SetHeader("X-Count", strconv.FormatInt(count, 10)).
		SetHeader("X-Date", date).
		Get(c.urlPrefix + urlFetchVolumeRankByDate)
	if err != nil {
		return data, err
	} else if resp.StatusCode() != http.StatusOK {
		return data, errors.New("FetchVolumeRankByDate API Fail")
	}
	body := VolumeRankArrProto{}
	if err = proto.Unmarshal(resp.Body(), &body); err != nil {
		return data, err
	}
	return body.Data, err
}

// FetchKbarByDateRange FetchKbarByDateRange
func (c *TradeAgent) FetchKbarByDateRange(stockNum string, start, end time.Time) (data []*KbarProto, err error) {
	stockAndDateArr := FetchKbarBody{
		StockNum:  stockNum,
		StartDate: start.Format(shortTimeLayout),
		EndDate:   end.Format(shortTimeLayout),
	}
	resp, err := c.restAgent.R().
		SetBody(stockAndDateArr).
		Post(c.urlPrefix + urlFetchKbarByDateRange)
	if err != nil {
		return data, err
	} else if resp.StatusCode() != http.StatusOK {
		return data, errors.New("FetchKbarByDateRange API Fail")
	}
	res := KbarArrProto{}
	if err = proto.Unmarshal(resp.Body(), &res); err != nil {
		return data, err
	}
	return res.Data, err
}

// FetchTSE001KbarByDate FetchTSE001KbarByDate
func (c *TradeAgent) FetchTSE001KbarByDate(date time.Time) (err error) {
	stockAndDateArr := FetchKbarBody{
		StartDate: date.Format(shortTimeLayout),
		EndDate:   date.Format(shortTimeLayout),
	}
	resp, err := c.restAgent.R().
		SetBody(stockAndDateArr).
		Post(c.urlPrefix + urlFetchTSE001KbarByDate)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("FetchTSEKbarByDateRange API Fail")
	}
	res := KbarArrProto{}
	if err = proto.Unmarshal(resp.Body(), &res); err != nil {
		return err
	}
	return err
}

// FetchEntireTickByStockAndDate FetchEntireTickByStockAndDate
func (c *TradeAgent) FetchEntireTickByStockAndDate(stockNum, date string) (data []*EntireTickProto, err error) {
	stockAndDate := FetchBody{
		StockNum: stockNum,
		Date:     date,
	}
	var resp *resty.Response
	resp, err = c.restAgent.R().
		SetBody(stockAndDate).
		Post(c.urlPrefix + urlFetchEntireTickByStockAndDate)
	if err != nil {
		return data, err
	} else if resp.StatusCode() != http.StatusOK {
		return data, errors.New("FetchEntireTickByStockAndDate API Fail")
	}
	res := EntireTickArrProto{}
	if err = proto.Unmarshal(resp.Body(), &res); err != nil {
		return data, err
	}
	return res.Data, err
}

// FetchAllStockDetail FetchAllStockDetail
func (c *TradeAgent) FetchAllStockDetail() (data []FetchStockBody, err error) {
	var resp *resty.Response
	resp, err = c.restAgent.R().
		SetResult(&[]FetchStockBody{}).
		Get(c.urlPrefix + urlFetchAllStockDetail)
	if err != nil {
		return data, err
	} else if resp.StatusCode() != http.StatusOK {
		return data, errors.New("FetchAllStockDetail API Fail")
	}
	return *resp.Result().(*[]FetchStockBody), err
}

// SubStreamTick SubStreamTick
func (c *TradeAgent) SubStreamTick(stockArr []string) (err error) {
	stocks := SubscribeBody{
		StockNumArr: stockArr,
	}
	var resp *resty.Response
	resp, err = c.restAgent.R().
		SetBody(stocks).
		SetResult(&OrderResponse{}).
		Post(c.urlPrefix + urlSubStreamTick)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("SubStreamTick API Fail")
	}
	res := *resp.Result().(*OrderResponse)
	if res.Status != StatusSuccuss {
		return errors.New(res.Status)
	}
	return err
}

// SubBidAsk SubBidAsk
func (c *TradeAgent) SubBidAsk(stockArr []string) (err error) {
	stocks := SubscribeBody{
		StockNumArr: stockArr,
	}
	var resp *resty.Response
	resp, err = c.restAgent.R().
		SetBody(stocks).
		SetResult(&OrderResponse{}).
		Post(c.urlPrefix + urlSubBidAsk)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("SubBidAsk API Fail")
	}
	if status := resp.Result().(*OrderResponse).Status; status != StatusSuccuss {
		return errors.New(status)
	}
	return err
}

// UnSubscribeAllByType UnSubscribeAllByType
func (c *TradeAgent) UnSubscribeAllByType(dataType TickType) (err error) {
	var url string
	switch {
	case dataType == StreamType:
		url = urlUnSubscribeAllStream
	case dataType == BidAsk:
		url = urlUnSubscribeAllBidAsk
	}
	var resp *resty.Response
	resp, err = c.restAgent.R().
		SetResult(&OrderResponse{}).
		Get(c.urlPrefix + url)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("UnSubscribeAllByType API Fail")
	}
	if status := resp.Result().(*OrderResponse).Status; status != StatusSuccuss {
		return errors.New(status)
	}
	return err
}
