// Package simulate package simulate
package simulate

import (
	"sort"

	"gitlab.tocraw.com/root/toc_trader/pkg/models/simulationcond"
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

// Update Update
func Update(result *Result, db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&result).Error; err != nil {
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

// GetBestForwardSimulateResult GetBestForwardSimulateResult
func GetBestForwardSimulateResult(db *gorm.DB) (cond Result, err error) {
	var beforeSort []Result
	err = db.Preload("Cond").
		Where("positive_days = total_days").
		Where("trade_count != positive_days").
		Where("forward_balance != 0").
		Order("(balance-total_loss)/trade_count desc").
		Find(&beforeSort).Error
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
			return afterSort[i].Cond.RsiHigh > afterSort[j].Cond.RsiHigh
		})
	}
	return afterSort[0], err
}

// GetBestReverseSimulateResult GetBestReverseSimulateResult
func GetBestReverseSimulateResult(db *gorm.DB) (cond Result, err error) {
	var beforeSort []Result
	err = db.Preload("Cond").
		Where("positive_days = total_days").
		Where("trade_count != positive_days").
		Where("reverse_balance != 0").
		Order("(balance-total_loss)/trade_count desc").
		Find(&beforeSort).Error
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
			return afterSort[i].Cond.RsiLow < afterSort[j].Cond.RsiLow
		})
	}
	return afterSort[0], err
}

// GetBestForwardCond GetBestForwardCond
func GetBestForwardCond(db *gorm.DB) (cond simulationcond.AnalyzeCondition, err error) {
	var beforeSort []Result
	err = db.Preload("Cond").Where("is_best_forward = true").Find(&beforeSort).Error
	if err != nil {
		return cond, err
	}
	if len(beforeSort) == 0 {
		return cond, err
	}
	return beforeSort[0].Cond, err
}

// GetBestReverseCond GetBestReverseCond
func GetBestReverseCond(db *gorm.DB) (cond simulationcond.AnalyzeCondition, err error) {
	var beforeSort []Result
	err = db.Preload("Cond").Where("is_best_reverse = true").Find(&beforeSort).Error
	if err != nil {
		return cond, err
	}
	if len(beforeSort) == 0 {
		return cond, err
	}
	return beforeSort[0].Cond, err
}
