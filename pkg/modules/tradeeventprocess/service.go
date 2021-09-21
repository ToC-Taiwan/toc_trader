// Package tradeeventprocess package tradeeventprocess
package tradeeventprocess

import (
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/tradeevent"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
)

// TradeEventSaver TradeEventSaver
func TradeEventSaver(record tradeevent.EventResponse) (err error) {
	if err = tradeevent.Insert(record, global.GlobalDB); err != nil {
		logger.Logger.Error(err)
		return err
	}
	return err
}

// CleanEvent CleanEvent
func CleanEvent() error {
	err := tradeevent.DeleteAll(global.GlobalDB)
	return err
}
