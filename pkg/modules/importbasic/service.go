// Package importbasic package importbasic
package importbasic

import (
	"errors"
	"net/http"
	"runtime/debug"
	"time"

	"gitlab.tocraw.com/root/toc_trader/external/sinopacsrv"
	"gitlab.tocraw.com/root/toc_trader/internal/database"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/internal/restful"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/holiday"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/stock"
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
			logger.GetLogger().Error(err.Error() + "\n" + string(debug.Stack()))
		}
	}()
	// Update basic first
	err = AskSinoSRVUpdateBasic()
	if err != nil {
		panic(err)
	}
	allStockNumInDB, err := stock.GetAllStockNum(database.GetAgent())
	if err != nil {
		panic(err)
	}
	existMap := make(map[string]bool)
	for _, v := range allStockNumInDB {
		existMap[v] = true
	}
	// Get stock detail from Sinopac SRV
	resp, err := restful.GetClient().R().
		SetResult(&[]sinopacsrv.FetchStockBody{}).
		Get("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/basic/importstock")
	if err != nil {
		panic(err)
	} else if resp.StatusCode() != http.StatusOK {
		panic("ImportAllStock api fail")
	}
	res := *resp.Result().(*[]sinopacsrv.FetchStockBody)
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
	if err := stock.InsertMultiRecord(insertArr, database.GetAgent()); err != nil {
		panic(err)
	}
	logger.GetLogger().WithFields(map[string]interface{}{
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

// AskSinoPyUpdateBasic AskSinoPyUpdateBasic
func AskSinoSRVUpdateBasic() (err error) {
	resp, err := restful.GetClient().R().
		SetResult(&sinopacsrv.OrderResponse{}).
		Get("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/basic/update-basic")
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("AskSinoPyUpdateBasic api fail")
	}
	res := *resp.Result().(*sinopacsrv.OrderResponse)
	if res.Status != sinopacsrv.StatusSuccuss {
		err = errors.New("sinopac srv update basic fail")
	}
	return err
}
