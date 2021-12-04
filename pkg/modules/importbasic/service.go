// Package importbasic package importbasic
package importbasic

import (
	"time"

	"gitlab.tocraw.com/root/toc_trader/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/database"
	"gitlab.tocraw.com/root/toc_trader/pkg/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/holiday"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/stock"
	"gitlab.tocraw.com/root/toc_trader/pkg/sinopacapi"
)

// AllStockDetailMap AllStockDetailMap
var AllStockDetailMap stock.MutexStruct

// ImportAllStock ImportAllStock
func ImportAllStock() (err error) {
	var allStockNumInDB []string
	allStockNumInDB, err = stock.GetAllStockNum(database.GetAgent())
	if err != nil {
		return err
	}

	existMap := make(map[string]bool)
	for _, v := range allStockNumInDB {
		existMap[v] = true
	}

	var res []sinopacapi.FetchStockBody
	res, err = sinopacapi.GetAgent().FetchAllStockDetail()
	if err != nil {
		return err
	}
	var importStock, already int64
	var insertArr []stock.Stock
	for _, v := range res {
		var dayTradeBool bool
		if v.DayTrade == "Yes" {
			dayTradeBool = true
		} else {
			continue
		}
		global.AllStockNameMap.Set(v.Code, v.Name)
		stock := stock.Stock{
			StockNum:  v.Code,
			StockName: v.Name,
			StockType: v.Exchange,
			DayTrade:  dayTradeBool,
			LastClose: v.Close,
			Category:  v.Category,
		}
		// Save detail in map
		AllStockDetailMap.Set(stock)
		if existMap[v.Code] {
			already++
		} else {
			insertArr = append(insertArr, stock)
			importStock++
		}
	}
	if err = stock.InsertMultiRecord(insertArr, database.GetAgent()); err != nil {
		return err
	}
	logger.GetLogger().WithFields(map[string]interface{}{
		"Imported":     importStock,
		"AlreadyExist": already,
	}).Info("Import Stock Status")
	return err
}

// ImportHoliday ImportHoliday
func ImportHoliday() (err error) {
	holidays := []string{
		"2021-01-01", "2021-02-08", "2021-02-09",
		"2021-02-10", "2021-02-11", "2021-02-12",
		"2021-02-15", "2021-02-16", "2021-03-01",
		"2021-04-02", "2021-04-05", "2021-04-30",
		"2021-06-14", "2021-09-20", "2021-09-21",
		"2021-10-11", "2021-12-31",
	}
	var holidayUTCArr []int64
	weekendArr := GetAllWeekend()
	for _, h := range holidays {
		var holidayUTC time.Time
		holidayUTC, err = time.Parse(global.ShortTimeLayout, h)
		if err != nil {
			return err
		}
		holidayUTCArr = append(holidayUTCArr, holidayUTC.Unix())
		for i, v := range weekendArr {
			if holidayUTC.Unix() == v {
				tmp := weekendArr[i+1:]
				weekendArr = weekendArr[:i]
				weekendArr = append(weekendArr, tmp...)
			}
		}
	}
	holidayUTCArr = append(holidayUTCArr, weekendArr...)

	allHolidayInDB, err := holiday.GetAllHoliday(database.GetAgent())
	if err != nil {
		return err
	}
	existMap := make(map[int64]bool)
	for _, v := range allHolidayInDB {
		existMap[v.TimeStamp] = true
	}
	var insertArr []holiday.Holiday
	for _, v := range holidayUTCArr {
		if existMap[v] {
			continue
		} else {
			tmp := holiday.Holiday{
				TimeStamp: v,
			}
			insertArr = append(insertArr, tmp)
		}
	}
	if err = holiday.InsertMultiRecord(insertArr, database.GetAgent()); err != nil {
		return err
	}
	return err
}

// GetAllWeekend GetAllWeekend
func GetAllWeekend() (weekendArr []int64) {
	firstDay := time.Date(int(global.TradeYear), 1, 1, 0, 0, 0, 0, time.UTC)
	for {
		if firstDay.Year() > int(global.TradeYear) {
			break
		}
		if firstDay.Weekday() == time.Saturday || firstDay.Weekday() == time.Sunday {
			weekendArr = append(weekendArr, firstDay.Unix())
		}
		firstDay = firstDay.AddDate(0, 0, 1)
	}
	return weekendArr
}

// GetTradeDay GetTradeDay
func GetTradeDay() (tradeDay time.Time, err error) {
	var today time.Time
	if time.Now().Hour() >= 15 {
		today = time.Now().AddDate(0, 0, 1)
	} else {
		today = time.Now()
	}
	tradeDay, err = GetNextTradeDayTime(today)
	if err != nil {
		return tradeDay, err
	}
	return tradeDay, err
}

// GetLastNTradeDay GetLastNTradeDay
func GetLastNTradeDay(n int64) (lastTradeDayArr []time.Time, err error) {
	var thisTradeDay, tmp time.Time
	if thisTradeDay, err = GetTradeDay(); err != nil {
		return lastTradeDayArr, err
	}
	for {
		if len(lastTradeDayArr) == int(n) {
			break
		}
		tmp, err = GetLastTradeDayTime(thisTradeDay)
		if err != nil {
			lastTradeDayArr = []time.Time{}
			return lastTradeDayArr, err
		}
		lastTradeDayArr = append(lastTradeDayArr, tmp)
		thisTradeDay = tmp
	}
	return lastTradeDayArr, err
}

// GetLastNTradeDayByDate GetLastNTradeDayByDate
func GetLastNTradeDayByDate(tradeDate time.Time, n int64) (lastTradeDayArr []time.Time, err error) {
	var tmp time.Time
	for {
		if len(lastTradeDayArr) == int(n) {
			break
		}
		tmp, err = GetLastTradeDayTime(tradeDate)
		if err != nil {
			return lastTradeDayArr, err
		}
		lastTradeDayArr = append(lastTradeDayArr, tmp)
		tradeDate = tmp
	}
	return lastTradeDayArr, err
}

// GetNextTradeDayTime GetNextTradeDayTime
func GetNextTradeDayTime(nowTime time.Time) (tradeDay time.Time, err error) {
	tmp := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 0, 0, 0, 0, time.UTC)
	exist, err := holiday.CheckIsHolidayByTimeStamp(tmp.Unix(), database.GetAgent())
	if err != nil {
		return tradeDay, err
	}
	if exist {
		nowTime = nowTime.AddDate(0, 0, 1)
		return GetNextTradeDayTime(nowTime)
	}
	return tmp, err
}

// GetLastTradeDayTime GetLastTradeDayTime
func GetLastTradeDayTime(tradeDay time.Time) (lastTradeDay time.Time, err error) {
	tmp := time.Date(tradeDay.Year(), tradeDay.Month(), tradeDay.Day(), 0, 0, 0, 0, time.UTC)
	exist, err := holiday.CheckIsHolidayByTimeStamp(tmp.AddDate(0, 0, -1).Unix(), database.GetAgent())
	if err != nil {
		return lastTradeDay, err
	}
	if exist {
		return GetLastTradeDayTime(tmp.AddDate(0, 0, -1))
	}
	return tmp.AddDate(0, 0, -1), err
}

// GetLastTradeDayByDate GetLastTradeDayByDate
func GetLastTradeDayByDate(tradeDay string) (lastTradeDay time.Time, err error) {
	dateUnix, err := time.Parse(global.ShortTimeLayout, tradeDay)
	if err != nil {
		return lastTradeDay, err
	}
	tmp := time.Date(dateUnix.Year(), dateUnix.Month(), dateUnix.Day(), 0, 0, 0, 0, time.UTC)
	exist, err := holiday.CheckIsHolidayByTimeStamp(tmp.AddDate(0, 0, -1).Unix(), database.GetAgent())
	if err != nil {
		return lastTradeDay, err
	}
	if exist {
		return GetLastTradeDayTime(tmp.AddDate(0, 0, -1))
	}
	return tmp.AddDate(0, 0, -1), err
}
