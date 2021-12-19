// Package choosetarget package choosetarget
package choosetarget

import (
	"sort"
	"sync"
	"time"

	"gitlab.tocraw.com/root/toc_trader/global"
	"gitlab.tocraw.com/root/toc_trader/init/sysparminit"
	"gitlab.tocraw.com/root/toc_trader/pkg/database"
	"gitlab.tocraw.com/root/toc_trader/pkg/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/entiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/snapshot"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/biasrate"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/fetchentiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/importbasic"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/subscribe"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/tradebot"
	"gitlab.tocraw.com/root/toc_trader/pkg/sinopacapi"
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
			if err = sinopacapi.GetAgent().RestartSinopacSRV(); err != nil && tradebot.BuyOrderMap.GetCount() != 0 && tradebot.SellFirstOrderMap.GetCount() != 0 {
				logger.GetLogger().Panic(err)
			}
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
	if err = biasrate.GetBiasRateByStockNumAndDate(*targetArr, global.TradeDay, 7); err != nil {
		logger.GetLogger().Error(err)
		return
	}
	subscribe.SubStockStreamTick(*targetArr)
}

// UnSubscribeAllFromSinopac UnSubscribeAllFromSinopac
func UnSubscribeAllFromSinopac() (err error) {
	// subscribe.UnSubscribeAll(subscribe.BidAsk)
	return subscribe.UnSubscribeStockAllByType(sinopacapi.StreamType)
}

// GetTopTarget GetTopTarget
func GetTopTarget(count int) (targetArr []string, err error) {
	var res []*sinopacapi.SnapShotProto
	res, err = sinopacapi.GetAgent().FetchAllSnapShot()
	if err != nil {
		return targetArr, err
	}
	var rank []*sinopacapi.SnapShotProto
	blackStockMap := sysparminit.GlobalSettings.GetBlackStockMap()
	blackCategoryMap := sysparminit.GlobalSettings.GetBlackCategoryMap()
	conditionArr := sysparminit.GlobalSettings.GetTargetCondArr()
	for _, condition := range conditionArr {
		for _, v := range res {
			if v.GetCode() == "001" {
				tmp := &snapshot.SnapShot{
					TS:              v.GetTs(),
					Code:            v.GetCode(),
					Exchange:        v.GetExchange(),
					Open:            v.GetOpen(),
					High:            v.GetHigh(),
					Low:             v.GetLow(),
					Close:           v.GetClose(),
					TickType:        v.GetTickType(),
					ChangePrice:     v.GetChangePrice(),
					ChangeRate:      v.GetChangeRate(),
					ChangeType:      v.GetChangeType(),
					AveragePrice:    v.GetAveragePrice(),
					Volume:          v.GetVolume(),
					TotalVolume:     v.GetTotalVolume(),
					Amount:          v.GetAmount(),
					TotalAmount:     v.GetTotalAmount(),
					YesterdayVolume: v.GetYesterdayVolume(),
					BuyPrice:        v.GetBuyPrice(),
					BuyVolume:       v.GetBuyVolume(),
					SellPrice:       v.GetSellPrice(),
					SellVolume:      v.GetSellVolume(),
					VolumeRatio:     v.GetVolumeRatio(),
				}
				TSEChannel <- tmp
				continue
			}
			if !importbasic.AllStockDetailMap.CheckIsDayTrade(v.GetCode()) {
				continue
			}
			if _, ok := blackStockMap[v.GetCode()]; ok {
				continue
			}
			if _, ok := blackCategoryMap[importbasic.AllStockDetailMap.GetCategory(v.GetCode())]; ok {
				continue
			}
			if count == -1 && v.GetTotalVolume() < condition.LimitVolumeLow {
				continue
			}
			if v.GetClose() > condition.LimitPriceLow && v.GetClose() < condition.LimitPriceHigh {
				rank = append(rank, v)
			}
		}
	}
	if len(rank) < 2 {
		return targetArr, err
	}
	sort.Slice(rank, func(i, j int) bool {
		return rank[i].GetTotalVolume() > rank[j].GetTotalVolume()
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
	stockCloseMap, needFetch, err := sinopacapi.GetAgent().FetchStockCloseMapByStockDateArr(stockNumArr, dateArr)
	if err != nil {
		return noCloseArr, err
	}
	for _, stock := range needFetch {
		for _, date := range dateArr {
			var tmpClose float64
			tmpClose, err = entiretick.GetLastCloseByDate(stock, date.Format(global.ShortTimeLayout), database.GetAgent())
			if err != nil {
				return noCloseArr, err
			}
			if tmpClose == 0 {
				var res []*entiretick.EntireTick
				res, err = fetchentiretick.FetchByDate(stock, date.Format(global.ShortTimeLayout))
				if err != nil {
					return noCloseArr, err
				}
				if len(res) == 0 {
					noCloseArr = append(noCloseArr, stock)
					logger.GetLogger().Warnf("%s cannot fetch %s close", stock, date.Format(global.ShortTimeLayout))
					break
				} else {
					fetchSaveLock.Lock()
					if err = entiretick.InsertMultiRecord(res, database.GetAgent()); err != nil {
						return noCloseArr, err
					}
					fetchSaveLock.Unlock()
					tmpClose = res[len(res)-1].Close
				}
			}
			global.StockCloseByDateMap.Set(date.Format(global.ShortTimeLayout), stock, tmpClose)
			logger.GetLogger().WithFields(map[string]interface{}{
				"Date":  date.Format(global.ShortTimeLayout),
				"Stock": stock,
				"Close": tmpClose,
			}).Infof("Fetch And Update Stock Close")
		}
	}
	for stock, dateMap := range stockCloseMap {
		for date, close := range dateMap {
			logger.GetLogger().WithFields(map[string]interface{}{
				"Date":  date,
				"Stock": stock,
				"Close": close,
			}).Infof("Update Stock Close")
			global.StockCloseByDateMap.Set(date, stock, close)
		}
	}
	return noCloseArr, err
}

// TSE001 TSE001
var TSE001 *snapshot.SnapShot

// TSEProcess TSEProcess
func TSEProcess() {
	lastClose, err := sinopacapi.GetAgent().FetchTSE001CloseByDate(global.LastTradeDay)
	if err != nil {
		logger.GetLogger().Error(err)
	}
	logger.GetLogger().Warnf("LastTradeDay %s TSE001 last close is %.2f", global.LastTradeDay.Format(global.ShortTimeLayout), lastClose)
	TSE001 = &snapshot.SnapShot{}
	for {
		tse := <-TSEChannel
		TSE001 = tse
	}
}

// GetVolumeRankByDate GetVolumeRankByDate
func GetVolumeRankByDate(date string, count int64) (rankArr []string, err error) {
	res, err := sinopacapi.GetAgent().FetchVolumeRankByDate(date, count)
	if err != nil {
		return rankArr, err
	}
	blackStockMap := sysparminit.GlobalSettings.GetBlackStockMap()
	blackCategoryMap := sysparminit.GlobalSettings.GetBlackCategoryMap()
	conditionArr := sysparminit.GlobalSettings.GetTargetCondArr()
	for _, v := range res {
		if !importbasic.AllStockDetailMap.CheckIsDayTrade(v.Code) {
			continue
		}
		if _, ok := blackStockMap[v.GetCode()]; ok {
			continue
		}
		if _, ok := blackCategoryMap[importbasic.AllStockDetailMap.GetCategory(v.GetCode())]; ok {
			continue
		}
		if v.GetTotalVolume() < conditionArr[0].LimitVolumeLow || v.GetTotalVolume() > conditionArr[0].LimitVolumeHigh {
			continue
		}
		if v.GetClose() > conditionArr[0].LimitPriceLow && v.GetClose() < conditionArr[0].LimitPriceHigh {
			rankArr = append(rankArr, v.GetCode())
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
		if newTargetArr, err := GetTopTarget(3); err != nil {
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
