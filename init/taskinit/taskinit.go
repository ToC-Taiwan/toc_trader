// Package taskinit init all task
package taskinit

import (
	"gitlab.tocraw.com/root/toc_trader/init/sysparminit"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/tasks/healthchecktask"
	"gitlab.tocraw.com/root/toc_trader/pkg/tasks/tradeeventtask"

	"github.com/robfig/cron"
)

func init() {
	c := cron.New()
	err := c.AddFunc(sysparminit.GlobalSettings.GetCleanEventCron(), func() {
		tradeeventtask.Run()
	})
	if err != nil {
		logger.GetLogger().Panic(err)
	}

	err = c.AddFunc(sysparminit.GlobalSettings.GetRestartSinopacCron(), func() {
		healthchecktask.Run()
	})
	if err != nil {
		logger.GetLogger().Panic(err)
	}
	c.Start()
}
