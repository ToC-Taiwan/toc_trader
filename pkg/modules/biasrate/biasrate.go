// Package biasrate package biasrate
package biasrate

import (
	"time"

	"gitlab.tocraw.com/root/toc_trader/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/importbasic"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tickanalyze"
	"gitlab.tocraw.com/root/toc_trader/pkg/sinopacapi"
	"gitlab.tocraw.com/root/toc_trader/pkg/utils"
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
	lastCloseMap, err := sinopacapi.GetAgent().FetchLastCloseByStockArrDateArr(stockNumArr, dateArr)
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
		biasRate := utils.Round(100*(closeArr[len(closeArr)-1]-ma)/ma, 2)
		result[stock] = biasRate
	}
	return result, err
}
