// Package tradeeventprocess package tradeeventprocess
package tradeeventprocess

import (
	"gitlab.tocraw.com/root/toc_trader/internal/database"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/tradeevent"
)

// TradeEventSaver TradeEventSaver
func TradeEventSaver(record tradeevent.EventResponse) (err error) {
	if err = tradeevent.Insert(record, database.GetAgent()); err != nil {
		logger.GetLogger().Error(err)
		return err
	}
	return err
}

// CleanEvent CleanEvent
func CleanEvent() error {
	err := tradeevent.DeleteAll(database.GetAgent())
	return err
}
