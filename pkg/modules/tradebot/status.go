// Package tradebot package tradebot
package tradebot

import (
	"errors"
	"time"

	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/pyresponse"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/traderecord"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
)

// CheckOrderStatusLoop CheckOrderStatusLoop
func CheckOrderStatusLoop() {
	go ShowStatus()
	tick := time.Tick(1*time.Second + 500*time.Millisecond)
	var initQuota bool
	for range tick {
		if err := FetchOrderStatus(); err != nil {
			logger.Logger.Error(err)
			continue
		}
		if !initQuota {
			if err := InitStartUpQuota(); err != nil {
				panic(err)
			}
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

// ShowStatus ShowStatus
func ShowStatus() {
	tick := time.Tick(60 * time.Second)
	for range tick {
		if lastTradeTime.IsZero() {
			lastTradeTime = time.Date(global.TradeDay.Year(), global.TradeDay.Month(), global.TradeDay.Day(), 13, 0, 0, 0, time.Local)
			logger.Logger.Infof("LastTradeTime is %s", lastTradeTime)
		}
		if time.Now().After(lastTradeTime) {
			global.EnableBuy = false
			logger.Logger.Warn("Trun enable buy off")
		}
		if FilledBuyOrderMap.GetCount() == FilledSellOrderMap.GetCount() && FilledSellOrderMap.GetCount() != 0 {
			balance := FilledSellOrderMap.GetTotalSellCost() - FilledBuyOrderMap.GetTotalBuyCost()
			logger.Logger.WithFields(map[string]interface{}{
				"Balance": balance,
			}).Info("Balance Status")
		}
		if BuyOrderMap.GetCount() != 0 {
			logger.Logger.WithFields(map[string]interface{}{
				"Current":       BuyOrderMap.GetCount(),
				"Maximum":       global.MeanTimeTradeStockNum,
				"TradeQuota":    TradeQuota,
				"LastTradeTime": lastTradeTime.Format(global.LongTimeLayout),
			}).Info("Current Trade Status")
		}
	}
}

func initBalance(orders []traderecord.TradeRecord) {
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
			logger.Logger.WithFields(map[string]interface{}{
				"StockNum": record.StockNum,
				"Name":     record.StockName,
				"Quantity": record.Quantity,
				"Price":    record.Price,
			}).Warn("Filled Buy Order")
		} else if val.Action == 2 && val.Status == 6 && !FilledSellOrderMap.CheckStockExist(val.StockNum) {
			FilledSellOrderMap.Set(record)
			logger.Logger.WithFields(map[string]interface{}{
				"StockNum": record.StockNum,
				"Name":     record.StockName,
				"Quantity": record.Quantity,
				"Price":    record.Price,
			}).Warn("Filled Sell Order")
		}
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
	if res.Status != pyresponse.SuccessStatus {
		return errors.New("FetchOrderStatus fail")
	}
	return err
}
