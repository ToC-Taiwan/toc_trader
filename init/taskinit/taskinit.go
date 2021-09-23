// Package taskinit init all task
package taskinit

import (
	"gitlab.tocraw.com/root/toc_trader/init/sysparminit"
	"gitlab.tocraw.com/root/toc_trader/pkg/tasks/tradeeventprocess"

	"github.com/robfig/cron"
)

func init() {
	c := cron.New()
	err := c.AddFunc(sysparminit.GlobalSettings.GetCleanEventCron(), func() {
		tradeeventprocess.Run()
	})
	if err != nil {
		panic(err)
	}

	err = c.AddFunc(sysparminit.GlobalSettings.GetRestartSinopacAndTocTraderCron(), func() {
		tradeeventprocess.Run()
	})
	if err != nil {
		panic(err)
	}
	c.Start()
}
