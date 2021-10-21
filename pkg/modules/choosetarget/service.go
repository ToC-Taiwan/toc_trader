// Package choosetarget package choosetarget
package choosetarget

import (
	"errors"
	"runtime/debug"
	"sort"
	"strconv"
	"sync"
	"time"

	"gitlab.tocraw.com/root/toc_trader/init/sysparminit"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/entiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/snapshot"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/stock"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/sysparm"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/targetstock"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/volumerank"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/fetchentiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/importbasic"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/subscribe"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tradebot"
	"gitlab.tocraw.com/root/toc_trader/tools/db"
	"gitlab.tocraw.com/root/toc_trader/tools/healthcheck"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
	"gitlab.tocraw.com/root/toc_trader/tools/rest"
	"google.golang.org/protobuf/proto"
)

// TSEChannel TSEChannel
var TSEChannel chan *snapshot.SnapShot

func init() {
	TSEChannel = make(chan *snapshot.SnapShot)
}

// SubscribeTarget SubscribeTarget
func SubscribeTarget(targetArr []string) {
	var errorTimes int
	// Update last 2 trade day close in map
	for {
		err := UpdateStockCloseMapByDate(targetArr, global.LastTradeDayArr)
		if errorTimes >= 5 {
			if err = healthcheck.FullRestart(); err != nil && tradebot.BuyOrderMap.GetCount() != 0 && tradebot.SellFirstOrderMap.GetCount() != 0 {
				logger.GetLogger().Fatal(err)
			}
			return
		}
		if err != nil {
			logger.GetLogger().Error(err)
			errorTimes++
		} else {
			break
		}
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
	resp, err := rest.GetClient().R().
		Get("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/basic/update/snapshot")
	if err != nil {
		return targetArr, err
	} else if resp.StatusCode() != 200 {
		return targetArr, errors.New("GetTopTarget api fail")
	}
	body := snapshot.SnapShotArrProto{}
	if err = proto.Unmarshal(resp.Body(), &body); err != nil {
		logger.GetLogger().Error(err)
		return targetArr, err
	}
	var rank []*snapshot.SnapShot
	blackStockMap := sysparminit.GlobalSettings.GetBlackStockMap()
	blackCategoryMap := sysparminit.GlobalSettings.GetBlackCategoryMap()
	conditionArr := sysparminit.GlobalSettings.GetTargetCondArr()
	for _, condition := range conditionArr {
		for _, v := range body.Data {
			if v.Code == "001" {
				TSEChannel <- v.ToSnapShotFromProto()
				continue
			}
			if !importbasic.AllStockDetailMap.CheckIsDayTrade(v.Code) {
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
	if len(rank) == 0 {
		return targetArr, err
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

// GetTargetFromStockList GetTargetFromStockList
func GetTargetFromStockList(conditionArr []sysparm.TargetCondArr) {
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
	var savedTarget []targetstock.Target
	if dbTarget, err := targetstock.GetTargetByTime(global.LastTradeDay, db.GetAgent()); err != nil {
		panic(err)
	} else if len(dbTarget) != 0 {
		for i, v := range dbTarget {
			global.TargetArr = append(global.TargetArr, v.Stock.StockNum)
			logger.GetLogger().WithFields(map[string]interface{}{
				"StockNum": v.Stock.StockNum,
				"Name":     v.Stock.StockName,
			}).Infof("Target %d", i+1)
		}
		return
	}
	blackStockMap := sysparminit.GlobalSettings.GetBlackStockMap()
	blackCategoryMap := sysparminit.GlobalSettings.GetBlackCategoryMap()
	if targets, err := stock.GetTargetByMultiLowHighVolume(conditionArr, db.GetAgent()); err != nil {
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
			logger.GetLogger().WithFields(map[string]interface{}{
				"StockNum": v.StockNum,
				"Name":     v.StockName,
			}).Infof("Target %d", i+1)
			savedTarget = append(savedTarget, targetstock.Target{
				LastTradeDay: global.LastTradeDay,
				Stock:        targets[i],
			})
		}
	}
	if err := targetstock.InsertMultiTarget(savedTarget, db.GetAgent()); err != nil {
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
			logger.GetLogger().Error(err.Error() + "\n" + string(debug.Stack()))
		}
	}()
	var lastUpdateTime time.Time
	if lastUpdateTime, err = stock.GetLastUpdatedTime(db.GetAgent()); err != nil {
		panic(err)
	} else if lastUpdateTime.After(global.LastTradeDay.Local().Add(7 * time.Hour)) {
		logger.GetLogger().Info("Volume and close is no need to update")
		return
	}
	resp, err := rest.GetClient().R().
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
	if err := stock.UpdateVolumeByStockNum(tmpMap, db.GetAgent()); err != nil {
		panic(err)
	}
	logger.GetLogger().Info("Volume and close update done")
}

var fetchSaveLock sync.RWMutex

// UpdateStockCloseMapByDate UpdateStockCloseMapByDate
func UpdateStockCloseMapByDate(stockNumArr []string, dateArr []time.Time) error {
	stockArr := FetchLastCountBody{
		StockNumArr: stockNumArr,
	}
	for _, date := range dateArr {
		logger.GetLogger().Infof("Update Stock Close on %s", date.Format(global.ShortTimeLayout))
		resp, err := rest.GetClient().R().
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
				tmpClose, err := entiretick.GetLastCloseByDate(val.Code, val.Date, db.GetAgent())
				if err != nil {
					return err
				}
				if tmpClose == 0 {
					res, err := fetchentiretick.FetchByDate(val.Code, val.Date)
					if err != nil {
						return err
					}
					if len(res) == 0 {
						logger.GetLogger().Errorf("%s cannot fetch %s close", val.Code, val.Date)
						continue
					} else {
						fetchSaveLock.Lock()
						if err := entiretick.InsertMultiRecord(res, db.GetAgent()); err != nil {
							return err
						}
						fetchSaveLock.Unlock()
						tmpClose = res[len(res)-1].Close
					}
				}
				global.StockCloseByDateMap.Set(val.Date, val.Code, tmpClose)
			}
		}
	}
	return nil
}

// TSE001 TSE001
var TSE001 *snapshot.SnapShot

// TSEProcess TSEProcess
func TSEProcess() {
	TSE001 = &snapshot.SnapShot{}
	for {
		tse := <-TSEChannel
		TSE001 = tse
	}
}

// GetTargetByVolumeRankByDate GetTargetByVolumeRankByDate
func GetTargetByVolumeRankByDate(date string, count int64) (rankArr []string, err error) {
	countStr := strconv.FormatInt(count, 10)
	resp, err := rest.GetClient().R().
		SetHeader("X-Count", countStr).
		SetHeader("X-Date", date).
		Get("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/trade/volumerank")
	if err != nil {
		return rankArr, err
	} else if resp.StatusCode() != 200 {
		return rankArr, errors.New("GetVolumeRankByDate api fail")
	}
	body := volumerank.VolumeRankArrProto{}
	if err = proto.Unmarshal(resp.Body(), &body); err != nil {
		logger.GetLogger().Error(err)
		return rankArr, err
	}
	blackStockMap := sysparminit.GlobalSettings.GetBlackStockMap()
	blackCategoryMap := sysparminit.GlobalSettings.GetBlackCategoryMap()
	conditionArr := sysparminit.GlobalSettings.GetTargetCondArr()
	for _, v := range body.Data {
		if !importbasic.AllStockDetailMap.CheckIsDayTrade(v.Code) {
			continue
		}
		if _, ok := blackStockMap[v.Code]; ok {
			continue
		}
		if _, ok := blackCategoryMap[importbasic.AllStockDetailMap.GetCategory(v.Code)]; ok {
			continue
		}
		if v.TotalVolume > conditionArr[0].LimitVolume && v.Close > conditionArr[0].LimitPriceLow && v.Close < conditionArr[0].LimitPriceHigh {
			rankArr = append(rankArr, v.Code)
		}
	}
	return rankArr, err
}
