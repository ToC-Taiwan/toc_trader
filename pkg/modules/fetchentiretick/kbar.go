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
	"gitlab.tocraw.com/root/toc_trader/pkg/models/kbar"
	"google.golang.org/protobuf/proto"
)

// FetchKbar FetchKbar
func FetchKbar(stockNumArr []string, start, end time.Time) {
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
	saveCh := make(chan *kbar.Kbar)
	go kbarSaver(saveCh)
	var wg sync.WaitGroup
	for _, v := range stockNumArr {
		wg.Add(1)
		stock := v
		go func(stockNum string) {
			defer wg.Done()
			var exist bool
			exist, err = kbar.CheckExistByStockAndDateRange(stockNum, start, end, database.GetAgent())
			if err != nil {
				logger.GetLogger().Panic(err)
			}
			if !exist {
				if err = kbar.DeleteByStockNum(stockNum, database.GetAgent()); err != nil {
					logger.GetLogger().Panic(err)
				}
				logger.GetLogger().WithFields(map[string]interface{}{
					"StockNum": stockNum,
					"From":     start.Format(global.ShortTimeLayout),
					"To":       end.Format(global.ShortTimeLayout),
				}).Info("Fetching Kbar")
				if err = FetchKbarByDateRange(stockNum, start, end, saveCh); err != nil {
					logger.GetLogger().Panic(err)
				}
			} else {
				logger.GetLogger().WithFields(map[string]interface{}{
					"StockNum": stockNum,
				}).Info("Kbar Already Exist")
			}
		}(stock)
	}
	wg.Wait()
	close(saveCh)
}

func kbarSaver(saveCh chan *kbar.Kbar) {
	var tmp []*kbar.Kbar
	for {
		kbarData, ok := <-saveCh
		if !ok {
			if err := kbar.InsertMultiRecord(tmp, database.GetAgent()); err != nil {
				logger.GetLogger().Error(err)
			}
			return
		}
		tmp = append(tmp, kbarData)
		if len(tmp) >= 2000 {
			if err := kbar.InsertMultiRecord(tmp, database.GetAgent()); err != nil {
				logger.GetLogger().Error(err)
				continue
			}
			tmp = []*kbar.Kbar{}
		}
	}
}

// FetchKbarByDateRange FetchKbarByDateRange
func FetchKbarByDateRange(stockNum string, start, end time.Time, saveCh chan *kbar.Kbar) (err error) {
	stockAndDateArr := FetchKbarBody{
		StockNum:  stockNum,
		StartDate: start.Format(global.ShortTimeLayout),
		EndDate:   end.Format(global.ShortTimeLayout),
	}
	resp, err := restful.GetClient().R().
		SetBody(stockAndDateArr).
		Post("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/history/kbar")
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("FetchKbarByDateRange api fail")
	}
	res := kbar.KbarArrProto{}
	if err = proto.Unmarshal(resp.Body(), &res); err != nil {
		return err
	}
	for _, v := range res.Data {
		saveCh <- v.ProtoToKbar(stockNum)
	}
	return err
}

// FetchTSEKbarByDate FetchTSEKbarByDate
func FetchTSEKbarByDate(date time.Time) (err error) {
	stockAndDateArr := FetchKbarBody{
		StartDate: date.Format(global.ShortTimeLayout),
		EndDate:   date.Format(global.ShortTimeLayout),
	}
	resp, err := restful.GetClient().R().
		SetBody(stockAndDateArr).
		Post("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/history/kbar/tse")
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("FetchTSEKbarByDateRange api fail")
	}
	res := kbar.KbarArrProto{}
	if err = proto.Unmarshal(resp.Body(), &res); err != nil {
		return err
	}
	return err
}
