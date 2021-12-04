// Package globalinit is init all global var
package globalinit

import (
	"gitlab.tocraw.com/root/toc_trader/pkg/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/importbasic"
)

func init() {
	if err := importbasic.ImportHoliday(); err != nil {
		logger.GetLogger().Panic(err)
	}
}
