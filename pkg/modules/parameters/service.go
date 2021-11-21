// Package parameters package parameters
package parameters

import (
	"gitlab.tocraw.com/root/toc_trader/init/sysparminit"
	"gitlab.tocraw.com/root/toc_trader/pkg/models/sysparm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// UpdateSysparm UpdateSysparm
func UpdateSysparm(key string, value interface{}) (err error) {
	db, err := gorm.Open(sqlite.Open(sysparminit.ConfigPath), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = sqlDB.Close(); err != nil {
			panic(err)
		}
	}()
	if err = sysparm.UpdateSysparm(key, value, db); err != nil {
		panic(err)
	}
	return err
}
