// Package importbasic package importbasic
package importbasic

import (
	"errors"
	"runtime/debug"
	"time"

	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/holiday"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/pyresponse"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/stock"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
)

// AllStockDetailMap AllStockDetailMap
var AllStockDetailMap stock.MutexStruct

// ImportAllStock ImportAllStock
func ImportAllStock() {
	var err error
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
			logger.Logger.Error(err.Error() + "\n" + string(debug.Stack()))
		}
	}()
	// Update basic first
	err = AskSinoPyUpdateBasic()
	if err != nil {
		panic(err)
	}
	allStockNumInDB, err := stock.GetAllStockNum(global.GlobalDB)
	if err != nil {
		panic(err)
	}
	existMap := make(map[string]bool)
	for _, v := range allStockNumInDB {
		existMap[v] = true
	}
	// Get stock detail from Sinopac SRV
	resp, err := global.RestyClient.R().
		SetResult(&[]FetchStockBody{}).
		Get("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/basic/importstock")
	if err != nil {
		panic(err)
	} else if resp.StatusCode() != 200 {
		panic("ImportAllStock api fail")
	}
	res := *resp.Result().(*[]FetchStockBody)
	var importStock, already int64
	var insertArr []stock.Stock
	for _, v := range res {
		stock := v.ToStock()
		// Save detail in map
		AllStockDetailMap.Set(stock)
		if existMap[v.Code] {
			already++
		} else {
			insertArr = append(insertArr, stock)
			importStock++
		}
	}
	if err := stock.InsertMultiRecord(insertArr, global.GlobalDB); err != nil {
		panic(err)
	}
	logger.Logger.WithFields(map[string]interface{}{
		"Imported":     importStock,
		"AlreadyExist": already,
	}).Info("Import Stock Status")
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

	allHolidayInDB, err := holiday.GetAllHoliday(global.GlobalDB)
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
	if err = holiday.InsertMultiRecord(insertArr, global.GlobalDB); err != nil {
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

// GetTradeDayTime GetTradeDayTime
func GetTradeDayTime(today time.Time) (tradeDay time.Time, err error) {
	tmp := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	exist, err := holiday.CheckIsHolidayByTimeStamp(tmp.Unix(), global.GlobalDB)
	if err != nil {
		return tradeDay, err
	}
	if exist {
		today = today.AddDate(0, 0, 1)
		return GetTradeDayTime(today)
	}
	return tmp, err
}

// GetLastTradeDayTime GetLastTradeDayTime
func GetLastTradeDayTime(tradeDay time.Time) (lastTradeDay time.Time, err error) {
	tmp := time.Date(tradeDay.Year(), tradeDay.Month(), tradeDay.Day(), 0, 0, 0, 0, time.UTC)
	exist, err := holiday.CheckIsHolidayByTimeStamp(tmp.AddDate(0, 0, -1).Unix(), global.GlobalDB)
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
	exist, err := holiday.CheckIsHolidayByTimeStamp(tmp.AddDate(0, 0, -1).Unix(), global.GlobalDB)
	if err != nil {
		return lastTradeDay, err
	}
	if exist {
		return GetLastTradeDayTime(tmp.AddDate(0, 0, -1))
	}
	return tmp.AddDate(0, 0, -1), err
}

// AskSinoPyUpdateBasic AskSinoPyUpdateBasic
func AskSinoPyUpdateBasic() (err error) {
	resp, err := global.RestyClient.R().
		SetResult(&pyresponse.PyServerResponse{}).
		Get("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/basic/update-basic")
	if err != nil {
		return err
	} else if resp.StatusCode() != 200 {
		return errors.New("AskSinoPyUpdateBasic api fail")
	}
	res := *resp.Result().(*pyresponse.PyServerResponse)
	if res.Status != pyresponse.SuccessStatus {
		err = errors.New("sinopac srv update basic fail")
	}
	return err
}
