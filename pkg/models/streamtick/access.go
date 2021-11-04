// Package streamtick package streamtick
package streamtick

import "gorm.io/gorm"

// Insert Insert
func Insert(tick StreamTick, db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&tick).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertMultiRecord InsertMultiRecord
func InsertMultiRecord(tickArr []*StreamTick, db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		batch := len(tickArr)
		if batch >= 2000 {
			batch = 2000
		}
		if err := tx.CreateInBatches(&tickArr, batch).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteAll DeleteAll
func DeleteAll(db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Not("id = 0").Unscoped().Delete(&StreamTick{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
