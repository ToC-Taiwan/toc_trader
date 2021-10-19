// Package dbinit init for db
package dbinit

import (
	"time"

	_ "github.com/lib/pq" // postgres driver for "database/sql"

	"database/sql"

	"gitlab.tocraw.com/root/toc_trader/init/sysparminit"
	"gitlab.tocraw.com/root/toc_trader/pkg/global"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/analyzeentiretick"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/analyzestreamtick"
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
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/parameters"
	"gitlab.tocraw.com/root/toc_trader/tools/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func init() {
	var err error
	initDataBase()
	dbLogger := gormlogger.New(logger.Logger, gormlogger.Config{
		SlowThreshold:             500 * time.Millisecond,
		Colorful:                  true,
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
	global.GlobalDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: dbLogger})
	if err != nil {
		panic(err)
	}

	err = global.GlobalDB.AutoMigrate(
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
}

func initDataBase() {
	db, err := sql.Open(
		"postgres",
		"user="+sysparminit.GlobalSettings.GetDBUser()+" password="+sysparminit.GlobalSettings.GetDBPass()+" host="+sysparminit.GlobalSettings.GetDBHost()+" port="+sysparminit.GlobalSettings.GetDBPort()+" sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			panic(err)
		}
	}()
	var exist bool
	statement := "SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = '" + sysparminit.GlobalSettings.GetDBName() + "')"
	err = db.QueryRow(statement).Scan(&exist)
	if err != nil {
		panic(err)
	}
	if !exist {
		_, err = db.Exec("CREATE DATABASE " + sysparminit.GlobalSettings.GetDBName() + ";")
		if err != nil {
			panic(err)
		}
	} else if sysparminit.GlobalSettings.GetResetParm() {
		_, err = db.Exec("DROP DATABASE " + sysparminit.GlobalSettings.GetDBName() + ";")
		if err != nil {
			panic(err)
		}
		_, err = db.Exec("CREATE DATABASE " + sysparminit.GlobalSettings.GetDBName() + ";")
		if err != nil {
			panic(err)
		}
		if err := parameters.UpdateSysparm("reset", 0); err != nil {
			panic(err)
		}
	}
}
