// Package balance package balance
package balance

import (
	"gorm.io/gorm"
)

// Insert Insert
func Insert(record Balance, db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertOrUpdate InsertOrUpdate
func InsertOrUpdate(record Balance, db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		var inDB Balance
		dbResult := db.Model(&Balance{}).Where("trade_day = ?", record.TradeDay).Find(&inDB)
		if dbResult.Error != nil {
			return dbResult.Error
		}
		if dbResult.RowsAffected == 0 {
			if err := tx.Create(&record).Error; err != nil {
				return err
			}
		} else {
			record.Model = inDB.Model
			if err := tx.Save(&record).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// InsertMultiRecord InsertMultiRecord
func InsertMultiRecord(recordArr []Balance, db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		batch := len(recordArr)
		if batch >= 2000 {
			batch = 2000
		}
		if err := tx.CreateInBatches(&recordArr, batch).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
