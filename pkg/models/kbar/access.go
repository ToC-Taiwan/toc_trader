// Package kbar package kbar
package kbar

import (
	"time"

	"gorm.io/gorm"
)

// Insert Insert
func Insert(kbar *Kbar, db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&kbar).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertMultiRecord InsertMultiRecord
func InsertMultiRecord(kbarArr []*Kbar, db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		batch := len(kbarArr)
		if batch >= 2000 {
			batch = 2000
		}
		if err := tx.CreateInBatches(&kbarArr, batch).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// CheckExistByStockAndDateRange CheckExistByStockAndDateRange
func CheckExistByStockAndDateRange(stockNum string, start, end time.Time, db *gorm.DB) (exist bool, err error) {
	var cnt1, cnt2 int64
	startTime := start.UnixNano()
	firstEnd := start.AddDate(0, 0, 1).UnixNano()
	endTime := end.UnixNano()
	secondEnd := end.AddDate(0, 0, 1).UnixNano()
	err = db.Model(&Kbar{}).Where("stock_num = ? AND timestamp >= ? AND timestamp < ?", stockNum, startTime, firstEnd).Count(&cnt1).Error
	if err != nil {
		return false, err
	}
	err = db.Model(&Kbar{}).Where("stock_num = ? AND timestamp >= ? AND timestamp < ?", stockNum, endTime, secondEnd).Count(&cnt2).Error
	if err != nil {
		return false, err
	}
	if cnt1 > 0 && cnt2 > 0 {
		return true, err
	}
	return exist, err
}

// DeleteByStockNum DeleteByStockNum
func DeleteByStockNum(stockNum string, db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("stock_num = ?", stockNum).Unscoped().Delete(&Kbar{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
