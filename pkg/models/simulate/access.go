// Package simulate package simulate
package simulate

import "gorm.io/gorm"

// Insert Insert
func Insert(result *Result, db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&result).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertMultiRecord InsertMultiRecord
func InsertMultiRecord(resultArr []Result, db *gorm.DB) error {
	batch := len(resultArr)
	if batch >= 2000 {
		batch = 2000
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(&resultArr, batch).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteAll DeleteAll
func DeleteAll(db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Not("id = 0").Unscoped().Delete(&Result{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
