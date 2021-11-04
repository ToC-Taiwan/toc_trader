// Package tickprocess package tickprocess
package tickprocess

import (
	"gitlab.tocraw.com/root/toc_trader/internal/database"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/streamtick"
)

// SaveStreamTicks SaveStreamTicks
func SaveStreamTicks(saveCh chan []*streamtick.StreamTick) {
	for {
		unSavedTicks := <-saveCh
		if len(unSavedTicks) != 0 {
			if err := streamtick.InsertMultiRecord(unSavedTicks, database.GetAgent()); err != nil {
				logger.GetLogger().Error(err)
				continue
			}
		}
	}
}

// DeleteAllStreamTicks DeleteAllStreamTicks
func DeleteAllStreamTicks() (err error) {
	if err = streamtick.DeleteAll(database.GetAgent()); err != nil {
		return err
	}
	return err
}
