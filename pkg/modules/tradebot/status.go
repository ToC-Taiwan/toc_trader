// Package tradebot package tradebot
package tradebot

import (
	"errors"
	"time"

	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/traderecord"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
)

// StatusFirstBack StatusFirstBack
var StatusFirstBack bool

// CheckOrderStatusLoop CheckOrderStatusLoop
func CheckOrderStatusLoop() {
	go showStatus()
	go tradeSwitch()
	tick := time.Tick(1*time.Second + 500*time.Millisecond)
	var initQuota bool
	for range tick {
		if err := FetchOrderStatus(); err != nil {
			logger.Logger.Error(err)
			continue
		}
		if !initQuota && StatusFirstBack {
			if err := InitStartUpQuota(); err != nil {
				panic(err)
			}
			logger.Logger.Warnf("Initial Quota: %d", TradeQuota)
			dbOrder, err := traderecord.GetAllorderByDayTime(global.TradeDay, global.GlobalDB)
			if err != nil {
				logger.Logger.Error(err)
				continue
			}
			initBalance(dbOrder)
			initQuota = true
		}
	}
}

// showStatus showStatus
func showStatus() {
	tick := time.Tick(60 * time.Second)
	for range tick {
		if isCurrentOrderAllFinished() && time.Now().Before(global.TradeDayEndTime.Add(2*time.Hour)) {
			balance := FilledSellOrderMap.GetTotalSellCost() + FilledSellFirstOrderMap.GetTotalSellCost() - FilledBuyLaterOrderMap.GetTotalBuyCost() - FilledBuyOrderMap.GetTotalBuyCost()
			back := FilledBuyOrderMap.GetTotalCostBack() + FilledSellOrderMap.GetTotalCostBack() + FilledSellFirstOrderMap.GetTotalCostBack() + FilledBuyLaterOrderMap.GetTotalCostBack()
			logger.Logger.WithFields(map[string]interface{}{
				"Current":         BuyOrderMap.GetCount(),
				"Maximum":         global.TradeSwitch.MeanTimeTradeStockNum,
				"TradeQuota":      TradeQuota,
				"OriginalBalance": balance,
				"Back":            back,
				"Real":            balance + back,
			}).Info("Current Trade Status")
			if balance < -1500 {
				global.TradeSwitch.Buy = false
				global.TradeSwitch.SellFirst = false
				logger.Logger.Warn("Enable buy and Sell first are all OFF because too...")
			}
		}
	}
}

func tradeSwitch() {
	tick := time.Tick(20 * time.Second)
	for range tick {
		if time.Now().After(global.TradeDayEndTime) && global.TradeSwitch.Buy {
			global.TradeSwitch.Buy = false
			global.TradeSwitch.SellFirst = false
			logger.Logger.Warn("Enable buy and Sell first are all OFF")
		}
	}
}

func isCurrentOrderAllFinished() bool {
	if BuyOrderMap.GetCount() != 0 || SellFirstOrderMap.GetCount() != 0 {
		return false
	}
	if FilledSellOrderMap.GetCount() == 0 && FilledBuyLaterOrderMap.GetCount() == 0 {
		return false
	}
	return FilledBuyOrderMap.GetCount() == FilledSellOrderMap.GetCount() && FilledSellFirstOrderMap.GetCount() == FilledBuyLaterOrderMap.GetCount()
}

func initBalance(orders []traderecord.TradeRecord) {
	var tmpBuyOrder, tmpSellOrder []traderecord.TradeRecord
	for _, val := range orders {
		record := traderecord.TradeRecord{
			StockNum:  val.StockNum,
			StockName: global.AllStockNameMap.GetName(val.StockNum),
			Action:    val.Action,
			Price:     val.Price,
			Quantity:  val.Quantity,
			Status:    val.Status,
			OrderID:   val.OrderID,
			TradeTime: time.Now(),
		}
		if val.Action == 1 && val.Status == 6 && !FilledBuyOrderMap.CheckStockExist(val.StockNum) {
			FilledBuyOrderMap.Set(record)
			tmpBuyOrder = append(tmpBuyOrder, record)
		} else if val.Action == 2 && val.Status == 6 && !FilledSellOrderMap.CheckStockExist(val.StockNum) {
			FilledSellOrderMap.Set(record)
			tmpSellOrder = append(tmpSellOrder, record)
		}
	}
	for _, v := range tmpBuyOrder {
		logger.Logger.WithFields(map[string]interface{}{
			"StockNum": v.StockNum,
			"Name":     v.StockName,
			"Quantity": v.Quantity,
			"Price":    v.Price,
		}).Warn("Filled Buy Order")
	}
	for _, v := range tmpSellOrder {
		logger.Logger.WithFields(map[string]interface{}{
			"StockNum": v.StockNum,
			"Name":     v.StockName,
			"Quantity": v.Quantity,
			"Price":    v.Price,
		}).Warn("Filled Sell Order")
	}
}

// FetchOrderStatus FetchOrderStatus
func FetchOrderStatus() (err error) {
	resp, err := global.RestyClient.R().
		SetResult(&traderecord.SinoStatusResponse{}).
		Get("http://" + global.PyServerHost + ":" + global.PyServerPort + "/pyapi/trade/status")
	if err != nil {
		return err
	} else if resp.StatusCode() != 200 {
		return errors.New("FetchOrderStatus api fail")
	}
	res := *resp.Result().(*traderecord.SinoStatusResponse)
	if res.Status != global.SuccessStatus {
		return errors.New("FetchOrderStatus fail")
	}
	return err
}
