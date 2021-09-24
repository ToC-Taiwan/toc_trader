// Package choosetarget package choosetarget
package choosetarget

import (
	"errors"
	"runtime/debug"
	"sort"
	"time"

	"gitlab.tocraw.com/root/toc_trader/init/sysparminit"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/snapshot"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/stock"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/sysparm"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/targetstock"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/fetchentiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/importbasic"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/subscribe"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
	"google.golang.org/protobuf/proto"
)

// TSEChannel TSEChannel
var TSEChannel chan *snapshot.SnapShot

func init() {
	TSEChannel = make(chan *snapshot.SnapShot)
}

// SubscribeTarget SubscribeTarget
func SubscribeTarget(targetArr []string) {
	// Update last 2 trade day close in map
	if err := UpdateStockCloseMapByDate(targetArr, global.LastTradeDayArr); err != nil {
		logger.Logger.Error(err)
	}
	// Subscribe all target stock and bidask and unsubscribefirst
	// subscribe.SubBidAsk(targetArr)
	subscribe.SubStreamTick(targetArr)
}

// UnSubscribeAll UnSubscribeAll
func UnSubscribeAll() {
	// subscribe.UnSubscribeAll(subscribe.BidAsk)
	subscribe.UnSubscribeAll(subscribe.StreamType)
}

// GetTopTarget GetTopTarget
func GetTopTarget(count int) (targetArr []string, err error) {
	go TSEProcess()
	resp, err := global.RestyClient.R().
		Get("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/basic/update/snapshot")
	if err != nil {
		return targetArr, err
	} else if resp.StatusCode() != 200 {
		return targetArr, errors.New("GetTopTarget api fail")
	}
	body := snapshot.SnapShotArrProto{}
	if err = proto.Unmarshal(resp.Body(), &body); err != nil {
		logger.Logger.Error(err)
		return targetArr, err
	}
	var rank []*snapshot.SnapShot
	blackStockMap := sysparminit.GlobalSettings.GetBlackStockMap()
	blackCategoryMap := sysparminit.GlobalSettings.GetBlackCategoryMap()
	conditionArr := sysparminit.GlobalSettings.GetTargetCondArr()
	for _, condition := range conditionArr {
		for _, v := range body.Data {
			if !importbasic.AllStockDetailMap.CheckIsDayTrade(v.Code) {
				continue
			}
			if v.Code == "001" {
				TSEChannel <- v.ToSnapShotFromProto()
				continue
			}
			if _, ok := blackStockMap[v.Code]; ok {
				continue
			}
			if _, ok := blackCategoryMap[importbasic.AllStockDetailMap.GetCategory(v.Code)]; ok {
				continue
			}
			if count == -1 && v.TotalVolume < condition.LimitVolume {
				continue
			}
			if v.Close > condition.LimitPriceLow && v.Close < condition.LimitPriceHigh {
				rank = append(rank, v.ToSnapShotFromProto())
			}
		}
	}
	sort.Slice(rank, func(i, j int) bool {
		return rank[i].TotalVolume > rank[j].TotalVolume
	})
	targetMap := make(map[string]bool)
	for _, target := range global.TargetArr {
		targetMap[target] = true
	}
	if count == -1 {
		if total := len(rank); total != 0 {
			count = total
		}
	}
	for _, stock := range rank[:count] {
		if _, ok := targetMap[stock.Code]; !ok {
			targetArr = append(targetArr, stock.Code)
		}
	}
	return targetArr, err
}

// GetTarget GetTarget
func GetTarget(conditionArr []sysparm.TargetCondArr) {
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
	var savedTarget []targetstock.Target
	if dbTarget, err := targetstock.GetTargetByTime(global.LastTradeDay, global.GlobalDB); err != nil {
		panic(err)
	} else if len(dbTarget) != 0 {
		for i, v := range dbTarget {
			global.TargetArr = append(global.TargetArr, v.Stock.StockNum)
			logger.Logger.WithFields(map[string]interface{}{
				"StockNum": v.Stock.StockNum,
				"Name":     v.Stock.StockName,
			}).Infof("Target %d", i+1)
		}
		return
	}
	blackStockMap := sysparminit.GlobalSettings.GetBlackStockMap()
	blackCategoryMap := sysparminit.GlobalSettings.GetBlackCategoryMap()
	if targets, err := stock.GetTargetByMultiLowHighVolume(conditionArr, global.GlobalDB); err != nil {
		panic(err)
	} else {
		for i, v := range targets {
			if _, ok := blackStockMap[v.StockNum]; ok {
				continue
			}
			if _, ok := blackCategoryMap[v.Category]; ok {
				continue
			}
			global.TargetArr = append(global.TargetArr, v.StockNum)
			logger.Logger.WithFields(map[string]interface{}{
				"StockNum": v.StockNum,
				"Name":     v.StockName,
			}).Infof("Target %d", i+1)
			savedTarget = append(savedTarget, targetstock.Target{
				LastTradeDay: global.LastTradeDay,
				Stock:        targets[i],
			})
		}
	}
	if err := targetstock.InsertMultiTarget(savedTarget, global.GlobalDB); err != nil {
		panic(err)
	}
}

// UpdateLastStockVolume UpdateLastStockVolume
func UpdateLastStockVolume() {
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
	var lastUpdateTime time.Time
	if lastUpdateTime, err = stock.GetLastUpdatedTime(global.GlobalDB); err != nil {
		panic(err)
	} else if lastUpdateTime.After(global.LastTradeDay.Local().Add(7 * time.Hour)) {
		logger.Logger.Info("Volume and close is no need to update")
		return
	}
	resp, err := global.RestyClient.R().
		Get("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/basic/update/snapshot")
	if err != nil {
		panic(err)
	} else if resp.StatusCode() != 200 {
		panic("UpdateLastStockVolume api fail")
	}
	tmpMap := make(map[string]struct {
		Close  float64
		Volume int64
	})
	res := snapshot.SnapShotArrProto{}
	if err := proto.Unmarshal(resp.Body(), &res); err != nil {
		panic(err)
	}
	for _, v := range res.Data {
		if v.Code == "001" {
			continue
		}
		tmpMap[v.Code] = struct {
			Close  float64
			Volume int64
		}{Close: v.GetClose(), Volume: v.GetTotalVolume()}
	}
	if err := stock.UpdateVolumeByStockNum(tmpMap, global.GlobalDB); err != nil {
		panic(err)
	}
	logger.Logger.Info("Volume and close update done")
}

// UpdateStockCloseMapByDate UpdateStockCloseMapByDate
func UpdateStockCloseMapByDate(stockNumArr []string, dateArr []time.Time) error {
	stockArr := FetchLastCountBody{
		StockNumArr: stockNumArr,
	}
	for _, date := range dateArr {
		logger.Logger.Infof("Update Stock Close on %s", date.Format(global.ShortTimeLayout))
		resp, err := global.RestyClient.R().
			SetHeader("X-Date", date.Format(global.ShortTimeLayout)).
			SetBody(stockArr).
			SetResult(&[]StockLastCount{}).
			Post("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/history/lastcount")
		if err != nil {
			return err
		} else if resp.StatusCode() != 200 {
			return errors.New("UpdateStockCloseMapByDate api fail")
		}
		stockLastCountArr := *resp.Result().(*[]StockLastCount)
		for _, val := range stockLastCountArr {
			if len(val.Close) != 0 {
				global.StockCloseByDateMap.Set(val.Date, val.Code, val.Close[0])
			} else {
				res, err := fetchentiretick.FetchByDate(val.Code, val.Date)
				if err != nil {
					return err
				}
				if len(res) == 0 {
					logger.Logger.Errorf("%s does not have %s close", val.Code, val.Date)
					continue
				}
				global.StockCloseByDateMap.Set(val.Date, val.Code, res[len(res)-1].Close)
			}
		}
	}
	return nil
}

// TSEProcess TSEProcess
func TSEProcess() {
	var tmp int64
	for {
		tse := <-TSEChannel
		if tmp != tse.TS {
			tmp = tse.TS
			logger.Logger.WithFields(map[string]interface{}{
				"Close":       tse.Close,
				"Open":        tse.Open,
				"High":        tse.High,
				"Low":         tse.Low,
				"ChangeRatio": tse.ChangeRate,
				"Diff":        tse.ChangePrice,
			}).Info("TSE")
		}
	}
}
