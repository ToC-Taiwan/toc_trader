// Package simulate package simulate
package simulate

import (
	"sort"

	"gorm.io/gorm"
)

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

// GetBestResult GetBestResult
func GetBestResult(db *gorm.DB) (cond Result, err error) {
	var beforeSort []Result
	err = db.Preload("Cond").Where("positive_days = total_days").
		Order("balance/trade_count desc").Find(&beforeSort).Error
	if err != nil {
		return cond, err
	}
	if len(beforeSort) == 0 {
		return cond, err
	}
	var afterSort []Result
	for _, v := range beforeSort {
		if v.Balance == beforeSort[0].Balance {
			afterSort = append(afterSort, v)
		} else {
			break
		}
	}
	if len(afterSort) > 1 {
		sort.Slice(afterSort, func(i, j int) bool {
			return afterSort[i].Cond.RsiHigh-afterSort[i].Cond.RsiLow > afterSort[j].Cond.RsiHigh-afterSort[j].Cond.RsiLow
		})
		sort.Slice(afterSort, func(i, j int) bool {
			return afterSort[i].Cond.RsiLow < afterSort[j].Cond.RsiLow
		})
	}
	return afterSort[0], err
}
