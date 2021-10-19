// Package sysparminit init sqlite parameter
package sysparminit

import (
	"os"
	"runtime"

	"gitlab.tocraw.com/root/toc_trader/pkg/models/sysparm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// GlobalSettings GlobalSettings
var GlobalSettings sysparm.GlobalSettingMap

func init() {
	GlobalSettings = make(map[string]string)
	db, err := gorm.Open(sqlite.Open("./configs/global.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			panic(err)
		}
	}()
	if err := db.AutoMigrate(&sysparm.Parameters{}); err != nil {
		panic(err)
	}

	var settings []sysparm.Parameters
	db.Model(&sysparm.Parameters{}).Find(&settings)
	for _, v := range settings {
		GlobalSettings[v.Key] = v.Value
	}
	for _, v := range sysparm.DefaultKey {
		if _, ok := GlobalSettings[v]; !ok {
			tmp := sysparm.Parameters{
				Key: v,
			}
			err := db.Transaction(func(tx *gorm.DB) error {
				if err := tx.Create(&tmp).Error; err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				panic(err)
			}
		}
	}

	if err := insertDefaultSetting(db); err != nil {
		panic(err)
	}
	db.Model(&sysparm.Parameters{}).Find(&settings)
	for _, v := range settings {
		GlobalSettings[v.Key] = v.Value
	}
}

func insertDefaultSetting(db *gorm.DB) (err error) {
	deployment := os.Getenv("DEPLOYMENT")
	if deployment == "docker" {
		sysparm.DefaultSetting["runmode"] = "release"
		sysparm.DefaultSetting["database"] = "tradebot"
		sysparm.DefaultSetting["dbhost"] = "172.20.10.10"
		sysparm.DefaultSetting["py_server_host"] = "sinopac-srv.tocraw.com"
	}
	if runtime.GOOS == "windows" {
		sysparm.DefaultSetting["dbhost"] = "172.20.10.10"
		sysparm.DefaultSetting["py_server_host"] = "sinopac-srv.tocraw.com"
	}
	var inDB []sysparm.Parameters
	db.Model(&sysparm.Parameters{}).Find(&inDB)
	for _, v := range inDB {
		key := v.Key
		if v.Value == "" {
			err = db.Transaction(func(tx *gorm.DB) error {
				if value, ok := sysparm.DefaultSetting[key]; ok {
					if err = tx.Model(&sysparm.Parameters{}).Where("key = ?", key).Update("value", value).Error; err != nil {
						return err
					}
				}
				return nil
			})
		}
	}
	return err
}
