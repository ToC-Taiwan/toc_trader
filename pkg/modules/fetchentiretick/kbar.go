// Package fetchentiretick package fetchentiretick
package fetchentiretick

import (
	"errors"
	"sync"
	"time"

	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/kbar"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
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
			exist, err := kbar.CheckExistByStockAndDateRange(stockNum, start, end, global.GlobalDB)
			if err != nil {
				panic(err)
			}
			if !exist {
				if err := kbar.DeleteByStockNum(stockNum, global.GlobalDB); err != nil {
					panic(err)
				}
				logger.GetLogger().Infof("Fetching %s's kbar from %s to %s", stockNum, start.Format(global.ShortTimeLayout), end.Format(global.ShortTimeLayout))
				if err := FetchKbarByDateRange(stockNum, start, end, saveCh); err != nil {
					panic(err)
				}
			} else {
				logger.GetLogger().Infof("%s Kbar from %s to %s already exist", stockNum, start.Format(global.ShortTimeLayout), end.Format(global.ShortTimeLayout))
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
			if err := kbar.InsertMultiRecord(tmp, global.GlobalDB); err != nil {
				logger.GetLogger().Error(err)
			}
			return
		}
		tmp = append(tmp, kbarData)
		if len(tmp) >= 2000 {
			if err := kbar.InsertMultiRecord(tmp, global.GlobalDB); err != nil {
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
	resp, err := global.RestyClient.R().
		SetBody(stockAndDateArr).
		Post("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/history/kbar")
	if err != nil {
		return err
	} else if resp.StatusCode() != 200 {
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
