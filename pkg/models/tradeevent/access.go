// Package tradeevent package tradeevent
package tradeevent

import "gorm.io/gorm"

// Insert Insert
func Insert(record EventResponse, db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteAll DeleteAll
func DeleteAll(db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Not("id = 0").Delete(&EventResponse{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
