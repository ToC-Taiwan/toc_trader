// Package fetchentiretick package fetchentiretick
package fetchentiretick

import (
	"errors"
	"net/http"
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
	saveCh := make(chan *kbar.Kbar)
	go kbarSaver(saveCh)
	var wg sync.WaitGroup
	for _, v := range stockNumArr {
		wg.Add(1)
		stock := v
		go func(stockNum string) {
			defer wg.Done()
			exist, err := kbar.CheckExistByStockAndDateRange(stockNum, start, end, database.GetAgent())
			if err != nil {
				panic(err)
			}
			if !exist {
				if err := kbar.DeleteByStockNum(stockNum, database.GetAgent()); err != nil {
					panic(err)
				}
				logger.GetLogger().WithFields(map[string]interface{}{
					"StockNum": stockNum,
					"From":     start.Format(global.ShortTimeLayout),
					"To":       end.Format(global.ShortTimeLayout),
				}).Info("Fetching Kbar")
				if err := FetchKbarByDateRange(stockNum, start, end, saveCh); err != nil {
					panic(err)
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
