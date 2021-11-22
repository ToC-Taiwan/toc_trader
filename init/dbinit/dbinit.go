// Package dbinit init for db
package dbinit

import (
	"database/sql"

	_ "github.com/lib/pq" // postgres driver for "database/sql"

	"gitlab.tocraw.com/root/toc_trader/init/sysparminit"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/modules/parameters"
)

func init() {
	db, err := sql.Open(
		"postgres",
		"user="+sysparminit.GlobalSettings.GetDBUser()+" password="+sysparminit.GlobalSettings.GetDBPass()+" host="+sysparminit.GlobalSettings.GetDBHost()+" port="+sysparminit.GlobalSettings.GetDBPort()+" sslmode=disable")
	if err != nil {
		logger.GetLogger().Panic(err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			logger.GetLogger().Panic(err)
		}
	}()
	var exist bool
	statement := "SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = '" + sysparminit.GlobalSettings.GetDBName() + "')"
	err = db.QueryRow(statement).Scan(&exist)
	if err != nil {
		logger.GetLogger().Panic(err)
	}
	if !exist {
		_, err = db.Exec("CREATE DATABASE " + sysparminit.GlobalSettings.GetDBName() + ";")
		if err != nil {
			logger.GetLogger().Panic(err)
		}
	} else if sysparminit.GlobalSettings.GetResetParm() {
		_, err = db.Exec("DROP DATABASE " + sysparminit.GlobalSettings.GetDBName() + ";")
		if err != nil {
			logger.GetLogger().Panic(err)
		}
		_, err = db.Exec("CREATE DATABASE " + sysparminit.GlobalSettings.GetDBName() + ";")
		if err != nil {
			logger.GetLogger().Panic(err)
		}
		if err := parameters.UpdateSysparm("reset", 0); err != nil {
			logger.GetLogger().Panic(err)
		}
	}
}
