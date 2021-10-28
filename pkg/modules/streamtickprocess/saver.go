// Package streamtickprocess package streamtickprocess
package streamtickprocess

import (
	"gitlab.tocraw.com/root/toc_trader/internal/db"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/streamtick"
)

// SaveStreamTicks SaveStreamTicks
func SaveStreamTicks(saveCh chan []*streamtick.StreamTick) {
	for {
		unSavedTicks := <-saveCh
		if len(unSavedTicks) != 0 {
			if err := streamtick.InsertMultiRecord(unSavedTicks, db.GetAgent()); err != nil {
				logger.GetLogger().Error(err)
				continue
			}
		}
	}
}
