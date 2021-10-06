// Package simulationcond package simulationcond
package simulationcond

import "gorm.io/gorm"

// Insert Insert
func Insert(data *AnalyzeCondition, db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&data).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertMultiRecord InsertMultiRecord
func InsertMultiRecord(dataArr []*AnalyzeCondition, db *gorm.DB) error {
	batch := len(dataArr)
	if batch >= 2000 {
		batch = 2000
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(&dataArr, batch).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteAll DeleteAll
func DeleteAll(db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Not("id = 0").Unscoped().Delete(&AnalyzeCondition{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
