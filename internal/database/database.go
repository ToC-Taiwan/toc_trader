// Package database package database
package database

import (
	"time"

	"gitlab.tocraw.com/root/toc_trader/init/sysparminit"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/analyzeentiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/analyzestreamtick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/balance"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/bidask"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/entiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/holiday"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/kbar"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulate"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/stock"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/streamtick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/targetstock"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/tradeevent"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/traderecord"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	gormlogger "gorm.io/gorm/logger"
)

// Agent Agent
var Agent *gorm.DB

// GetAgent GetAgent
func GetAgent() *gorm.DB {
	if Agent != nil {
		return Agent
	}
	var err error
	dbLogger := gormlogger.New(logger.GetLogger(), gormlogger.Config{
		SlowThreshold:             1000 * time.Millisecond,
		Colorful:                  false,
		IgnoreRecordNotFoundError: false,
		LogLevel:                  gormlogger.Warn,
	})
	dsn := "host=" + sysparminit.GlobalSettings.GetDBHost() +
		" user=" + sysparminit.GlobalSettings.GetDBUser() +
		" password=" + sysparminit.GlobalSettings.GetDBPass() +
		" dbname=" + sysparminit.GlobalSettings.GetDBName() +
		" port=" + sysparminit.GlobalSettings.GetDBPort() +
		" sslmode=disable" +
		" TimeZone=" + sysparminit.GlobalSettings.GetDBTimeZone()
	Agent, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: dbLogger, SkipDefaultTransaction: true})
	if err != nil {
		panic(err)
	}

	err = Agent.AutoMigrate(
		&balance.Balance{},
		&analyzeentiretick.AnalyzeEntireTick{},
		&analyzestreamtick.AnalyzeStreamTick{},
		&bidask.BidAsk{},
		&entiretick.EntireTick{},
		&holiday.Holiday{},
		&kbar.Kbar{},
		&simulate.Result{},
		&simulationcond.AnalyzeCondition{},
		&stock.Stock{},
		&streamtick.StreamTick{},
		&targetstock.Target{},
		&tradeevent.EventResponse{},
		&traderecord.TradeRecord{},
	)
	if err != nil {
		panic(err)
	}
	return Agent
}
