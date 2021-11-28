// Package choosetarget package choosetarget
package choosetarget

import (
	"errors"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"gitlab.tocraw.com/root/toc_trader/external/sinopacsrv"
	"gitlab.tocraw.com/root/toc_trader/init/sysparminit"
	"gitlab.tocraw.com/root/toc_trader/internal/database"
	"gitlab.tocraw.com/root/toc_trader/internal/healthcheck"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/internal/restful"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/entiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/snapshot"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/volumerank"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/biasrate"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/fetchentiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/importbasic"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/subscribe"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tradebot"
	"google.golang.org/protobuf/proto"
)

// TSEChannel TSEChannel
var TSEChannel chan *snapshot.SnapShot

func init() {
	TSEChannel = make(chan *snapshot.SnapShot)
}

// SubscribeTarget SubscribeTarget
func SubscribeTarget(targetArr *[]string) {
	var errorTimes int
	var noCloseArr []string
	var err error
	// Update last 2 trade day close in map
	for {
		noCloseArr, err = UpdateStockCloseMapByDate(*targetArr, global.LastTradeDayArr)
		if errorTimes >= 5 {
			if err = healthcheck.AskSinopacSRVRestart(); err != nil && tradebot.BuyOrderMap.GetCount() != 0 && tradebot.SellFirstOrderMap.GetCount() != 0 {
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
	// subscribe.SubBidAsk(targetArr)
	if len(noCloseArr) != 0 {
		tmp := make(map[string]bool)
		for _, v := range noCloseArr {
			tmp[v] = true
		}
		var subArr []string
		for _, k := range *targetArr {
			if _, ok := tmp[k]; !ok {
				subArr = append(subArr, k)
			}
		}
		*targetArr = subArr
	}
	// Fill BiasRate Map
	if err = biasrate.GetBiasRateByStockNumAndDate(*targetArr, global.TradeDay); err != nil {
		logger.GetLogger().Error(err)
		return
	}
	subscribe.SubStreamTick(*targetArr)
}

// UnSubscribeAll UnSubscribeAll
func UnSubscribeAll() {
	// subscribe.UnSubscribeAll(subscribe.BidAsk)
	subscribe.UnSubscribeAll(subscribe.StreamType)
}

// GetTopTarget GetTopTarget
func GetTopTarget(count int) (targetArr []string, err error) {
	resp, err := restful.GetClient().R().
		Get("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/basic/update/snapshot")
	if err != nil {
		return targetArr, err
	} else if resp.StatusCode() != http.StatusOK {
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

var fetchSaveLock sync.Mutex

// UpdateStockCloseMapByDate UpdateStockCloseMapByDate
func UpdateStockCloseMapByDate(stockNumArr []string, dateArr []time.Time) (noCloseArr []string, err error) {
	stockArr := FetchLastCountBody{
		StockNumArr: stockNumArr,
	}
	for _, date := range dateArr {
		logger.GetLogger().WithFields(map[string]interface{}{
			"Date": date.Format(global.ShortTimeLayout),
		}).Infof("Update Stock Close")
		var resp *resty.Response
		resp, err = restful.GetClient().R().
			SetHeader("X-Date", date.Format(global.ShortTimeLayout)).
			SetBody(stockArr).
			SetResult(&[]sinopacsrv.StockLastCount{}).
			Post("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/history/lastcount")
		if err != nil {
			return noCloseArr, err
		} else if resp.StatusCode() != http.StatusOK {
			return noCloseArr, errors.New("UpdateStockCloseMapByDate api fail")
		}
		stockLastCountArr := *resp.Result().(*[]sinopacsrv.StockLastCount)
		for _, val := range stockLastCountArr {
			if len(val.Close) != 0 {
				global.StockCloseByDateMap.Set(val.Date, val.Code, val.Close[0])
			} else {
				var tmpClose float64
				tmpClose, err = entiretick.GetLastCloseByDate(val.Code, val.Date, database.GetAgent())
				if err != nil {
					return noCloseArr, err
				}
				if tmpClose == 0 {
					var res []*entiretick.EntireTick
					res, err = fetchentiretick.FetchByDate(val.Code, val.Date)
					if err != nil {
						return noCloseArr, err
					}
					if len(res) == 0 {
						noCloseArr = append(noCloseArr, val.Code)
						logger.GetLogger().Warnf("%s cannot fetch %s close", val.Code, val.Date)
						continue
					} else {
						fetchSaveLock.Lock()
						if err = entiretick.InsertMultiRecord(res, database.GetAgent()); err != nil {
							return noCloseArr, err
						}
						fetchSaveLock.Unlock()
						tmpClose = res[len(res)-1].Close
					}
				}
				global.StockCloseByDateMap.Set(val.Date, val.Code, tmpClose)
			}
		}
	}
	return noCloseArr, err
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

// GetVolumeRankByDate GetVolumeRankByDate
func GetVolumeRankByDate(date string, count int64) (rankArr []string, err error) {
	resp, err := restful.GetClient().R().
		SetHeader("X-Count", strconv.FormatInt(count, 10)).
		SetHeader("X-Date", date).
		Get("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/trade/volumerank")
	if err != nil {
		return rankArr, err
	} else if resp.StatusCode() != http.StatusOK {
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

// AddTop10RankTarget AddTop10RankTarget
func AddTop10RankTarget() {
	tick := time.Tick(30 * time.Second)
	for range tick {
		if !tradebot.CheckIsOpenTime() {
			continue
		}
		var count int
		if newTargetArr, err := GetTopTarget(10); err != nil {
			logger.GetLogger().Error(err)
			continue
			// Start from 9:10 every 30 seconds
		} else if time.Now().After(global.TradeDay.Add(1*time.Hour + 10*time.Minute)) {
			count = len(newTargetArr)
			if count != 0 {
				SubscribeTarget(&newTargetArr)
				global.TargetArr = append(global.TargetArr, newTargetArr...)
			}
		}
		if count != 0 {
			logger.GetLogger().Infof("GetTopTarget %d", count)
		}
	}
}
