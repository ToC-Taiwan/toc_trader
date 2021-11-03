// Package fetchentiretick package fetchentiretick
package fetchentiretick

import (
	"errors"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"gitlab.tocraw.com/root/toc_trader/internal/database"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/internal/restful"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/entiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/importbasic"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tickprocess"
	"google.golang.org/protobuf/proto"
)

var wg sync.WaitGroup

// FetchEntireTick FetchEntireTick
func FetchEntireTick(stockNumArr []string, dateArr []time.Time, cond simulationcond.AnalyzeCondition) {
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
	saveCh := make(chan []*entiretick.EntireTick, len(stockNumArr))
	go tickprocess.SaveEntireTicks(saveCh)
	for _, d := range dateArr {
		for _, s := range stockNumArr {
			rows, err := entiretick.GetCntByStockAndDate(s, d.Format(global.ShortTimeLayout), database.GetAgent())
			if err != nil {
				panic(err)
			} else {
				if rows > 0 {
					logger.GetLogger().WithFields(map[string]interface{}{
						"Stock": s,
						"Date":  d.Format(global.ShortTimeLayout),
					}).Info("EntireTick Already Exist")
					continue
				} else {
					wg.Add(1)
					go GetAndSaveEntireTick(s, d.Format(global.ShortTimeLayout), cond, saveCh)
				}
			}
		}
		wg.Wait()
	}
	close(saveCh)
}

// GetAndSaveEntireTick GetAndSaveEntireTick
func GetAndSaveEntireTick(stockNum, date string, cond simulationcond.AnalyzeCondition, saveCh chan []*entiretick.EntireTick) {
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
	logger.GetLogger().WithFields(map[string]interface{}{
		"StockNum": stockNum,
		"Date":     date,
	}).Info("Fetching Entiretick")
	stockAndDate := FetchBody{
		StockNum: stockNum,
		Date:     date,
	}
	resp, err := restful.GetClient().R().
		SetBody(stockAndDate).
		Post("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/history/entiretick")
	if err != nil {
		panic(err)
	} else if resp.StatusCode() != http.StatusOK {
		panic("GetAndSaveEntireTick api fail")
	}
	res := entiretick.EntireTickArrProto{}
	if err = proto.Unmarshal(resp.Body(), &res); err != nil {
		panic(err)
	}
	ch := make(chan *entiretick.EntireTick, len(res.Data))
	lastTradeDay, err := importbasic.GetLastTradeDayByDate(date)
	if err != nil {
		panic(err)
	}
	var simulateMap tickprocess.AnalyzeEntireTickMap
	lastClose := global.StockCloseByDateMap.GetClose(stockNum, lastTradeDay.Format(global.ShortTimeLayout))
	if lastClose != 0 {
		go tickprocess.TickProcess(stockNum, lastClose, cond, ch, &wg, saveCh, false, &simulateMap)
	} else {
		logger.GetLogger().Warnf("%s has no %s's close", stockNum, date)
	}

	for _, tmpTick := range res.Data {
		tick, err := tmpTick.ProtoToEntireTick(stockNum)
		if err != nil {
			panic(err)
		}
		ch <- tick
	}
	close(ch)
}

// FetchByDate FetchByDate
func FetchByDate(stockNum, date string) (data []*entiretick.EntireTick, err error) {
	stockAndDateArr := FetchBody{
		StockNum: stockNum,
		Date:     date,
	}
	resp, err := restful.GetClient().R().
		SetBody(stockAndDateArr).
		Post("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/history/entiretick")
	if err != nil {
		return data, err
	} else if resp.StatusCode() != http.StatusOK {
		return data, errors.New("FetchByDate api fail")
	}
	res := entiretick.EntireTickArrProto{}
	if err = proto.Unmarshal(resp.Body(), &res); err != nil {
		return data, err
	}
	for _, v := range res.Data {
		var tick *entiretick.EntireTick
		tick, err = v.ProtoToEntireTick(stockNum)
		if err != nil {
			panic(err)
		}
		data = append(data, tick)
	}
	return data, err
}
