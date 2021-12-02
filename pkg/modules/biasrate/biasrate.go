// Package biasrate package biasrate
package biasrate

import (
	"errors"
	"net/http"
	"time"

	"gitlab.tocraw.com/root/toc_trader/external/sinopacsrv"
	"gitlab.tocraw.com/root/toc_trader/internal/common"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/internal/restful"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/importbasic"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tickanalyze"
)

// StockBiasRateMap StockBiasRateMap
var StockBiasRateMap MutexCache

// GetBiasRateByStockNumAndDate GetBiasRateByStockNumAndDate
func GetBiasRateByStockNumAndDate(stockNumArr []string, date time.Time, n int64) (err error) {
	logger.GetLogger().WithFields(map[string]interface{}{
		"TradeDay":   date.Format(global.ShortTimeLayout),
		"StockCount": len(stockNumArr),
		"Days":       n,
	}).Infof("Get %d Day BiasRate", n)
	tradeDayArr, err := importbasic.GetLastNTradeDayByDate(date, n)
	if err != nil {
		logger.GetLogger().Panic(err)
	}
	dateArr := []string{}
	for _, date := range tradeDayArr {
		dateArr = append(dateArr, date.Format(global.ShortTimeLayout))
	}
	lastCloseMap, err := FetchLastCloseByStockArrDateArr(stockNumArr, dateArr)
	if err != nil {
		return err
	}
	biasMap, err := GetBiasRate(lastCloseMap)
	if err != nil {
		return err
	}
	for stock, bias := range biasMap {
		logger.GetLogger().WithFields(map[string]interface{}{
			"Stock": stock,
			"Value": bias,
		}).Info("BiasRate")
		StockBiasRateMap.Set(stock, date.Format(global.ShortTimeLayout), bias)
	}
	return err
}

// GetBiasRate GetBiasRate
func GetBiasRate(lastCloseMap map[string][]float64) (result map[string]float64, err error) {
	result = make(map[string]float64)
	for stock, closeArr := range lastCloseMap {
		var ma float64
		ma, err = tickanalyze.GenerareMAByCount(closeArr, len(closeArr))
		if err != nil {
			return result, err
		}
		biasRate := common.Round(100*(closeArr[len(closeArr)-1]-ma)/ma, 2)
		result[stock] = biasRate
	}
	return result, err
}

// FetchLastCloseByStockArrDateArr FetchLastCloseByStockArrDateArr
func FetchLastCloseByStockArrDateArr(stockNumArr, dateArr []string) (stockLastCloseMap map[string][]float64, err error) {
	stockLastCloseMap = make(map[string][]float64)
	stockAndDateArr := FetchLastCloseBody{
		StockNumArr: stockNumArr,
		DateArr:     dateArr,
	}
	resp, err := restful.GetClient().R().
		SetBody(stockAndDateArr).
		SetResult(&[]sinopacsrv.LastCloseWithStockAndDate{}).
		Post("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/history/lastcount/multi-date")
	if err != nil {
		return stockLastCloseMap, err
	} else if resp.StatusCode() != http.StatusOK {
		return stockLastCloseMap, errors.New("FetchLastCloseByStockArrDateArr api fail")
	}
	res := *resp.Result().(*[]sinopacsrv.LastCloseWithStockAndDate)
	for _, v := range res {
		var tmp []float64
		for k := len(v.CloseArr) - 1; k >= 0; k-- {
			tmp = append(tmp, v.CloseArr[k].Close)
		}
		stockLastCloseMap[v.StockNum] = tmp
	}
	return stockLastCloseMap, err
}
