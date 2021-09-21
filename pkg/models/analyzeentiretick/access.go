// Package analyzeentiretick analyzeentiretick
package analyzeentiretick

import (
	"gorm.io/gorm"
)

// Insert Insert
func Insert(tick AnalyzeEntireTick, db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&tick).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertMultiRecord InsertMultiRecord
func InsertMultiRecord(tickArr []*AnalyzeEntireTick, db *gorm.DB) error {
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

// GetAllAnalyzeEntireTick GetAllAnalyzeEntireTick
func GetAllAnalyzeEntireTick(db *gorm.DB) (records []AnalyzeEntireTick, err error) {
	err = db.Model(&AnalyzeEntireTick{}).Order("timestamp asc").Find(&records).Error
	return records, err
}

// DeleteAll DeleteAll
func DeleteAll(db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Not("id = 0").Unscoped().Delete(&AnalyzeEntireTick{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
