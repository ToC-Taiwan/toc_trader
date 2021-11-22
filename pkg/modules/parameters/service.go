// Package parameters package parameters
package parameters

import (
	"database/sql"
	"errors"
	"runtime/debug"

	"gitlab.tocraw.com/root/toc_trader/init/sysparminit"
	"gitlab.tocraw.com/root/toc_trader/internal/logger"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/sysparm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// UpdateSysparm UpdateSysparm
func UpdateSysparm(key string, value interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
			logger.GetLogger().Error(err.Error() + "\n" + string(debug.Stack()))
		}
	}()
	var db *gorm.DB
	db, err = gorm.Open(sqlite.Open(sysparminit.ConfigPath), &gorm.Config{})
	if err != nil {
		logger.GetLogger().Panic(err)
	}
	var sqlDB *sql.DB
	sqlDB, err = db.DB()
	if err != nil {
		logger.GetLogger().Panic(err)
	}
	defer func() {
		if err = sqlDB.Close(); err != nil {
			logger.GetLogger().Panic(err)
		}
	}()
	if err = sysparm.UpdateSysparm(key, value, db); err != nil {
		logger.GetLogger().Panic(err)
	}
	return err
}
