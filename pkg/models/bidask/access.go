// Package bidask package bidask
package bidask

import "gorm.io/gorm"

// Insert Insert
func Insert(bidAsk BidAsk, db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&bidAsk).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertMultiRecord InsertMultiRecord
func InsertMultiRecord(bidAskArr []*BidAsk, db *gorm.DB) error {
	batch := len(bidAskArr)
	if batch >= 2000 {
		batch = 2000
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(&bidAskArr, batch).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
