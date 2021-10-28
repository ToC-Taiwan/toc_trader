// Package bidaskprocess package bidaskprocess
package bidaskprocess

import (
	"time"

	"gitlab.tocraw.com/root/toc_trader/internal/db"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/bidask"
)

// TmpBidAskMap TmpBidAskMap
var TmpBidAskMap BidAskMutexStruct

// SaveBidAsk SaveBidAsk
func SaveBidAsk(stockNum string) {
	tick := time.Tick(5 * time.Second)
	for range tick {
		if TmpBidAskMap.GetCountByStockNum(stockNum) == 0 {
			continue
		}
		tmpArr := TmpBidAskMap.GetArrByStockNum(stockNum)
		if err := bidask.InsertMultiRecord(tmpArr[:len(tmpArr)-1], db.GetAgent()); err != nil {
			logger.GetLogger().Error(err)
			continue
		}
		TmpBidAskMap.KeepLastOne(stockNum)
	}
}
