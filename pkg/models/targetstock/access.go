// Package targetstock package targetstock
package targetstock

import (
	"time"

	"gorm.io/gorm"
)

// Insert Insert
func Insert(target Target, db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&target).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertMultiRecord InsertMultiRecord
func InsertMultiRecord(targetArr []Target, db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		batch := len(targetArr)
		if batch >= 2000 {
			batch = 2000
		}
		if err := tx.CreateInBatches(&targetArr, batch).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// GetTargetByTime GetTargetByTime
func GetTargetByTime(tradeDayTime time.Time, db *gorm.DB) (data []Target, err error) {
	err = db.Preload("Stock").Where("trade_day = ?", tradeDayTime).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, err
}
