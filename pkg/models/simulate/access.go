// Package simulate package simulate
package simulate

import (
	"sort"
	"time"

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

// ClearIsBestForwardByTradeDay ClearIsBestForwardByTradeDay
func ClearIsBestForwardByTradeDay(tradeDay time.Time, db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("is_best_forward = true AND trade_day = ?", tradeDay).Unscoped().Delete(&Result{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// ClearIsBestReverseByTradeDay ClearIsBestReverseByTradeDay
func ClearIsBestReverseByTradeDay(tradeDay time.Time, db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("is_best_reverse = true AND trade_day = ?", tradeDay).Unscoped().Delete(&Result{}).Error; err != nil {
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

// DeleteAllNotBest DeleteAllNotBest
func DeleteAllNotBest(db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("is_best_forward = false AND is_best_reverse = false").Unscoped().Delete(&Result{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// GetBestForwardSimulateResultByTradeDay GetBestForwardSimulateResultByTradeDay
func GetBestForwardSimulateResultByTradeDay(tradeDay time.Time, db *gorm.DB) (cond Result, err error) {
	var beforeSort []Result
	err = db.Preload("Cond").
		Where("trade_day = ?", tradeDay).
		Where("positive_days = total_days").
		Where("trade_count != positive_days").
		Where("forward_balance != 0").
		Order("balance-total_loss desc").
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

// GetBestReverseSimulateResultByTradeDay GetBestReverseSimulateResultByTradeDay
func GetBestReverseSimulateResultByTradeDay(tradeDay time.Time, db *gorm.DB) (cond Result, err error) {
	var beforeSort []Result
	err = db.Preload("Cond").
		Where("trade_day = ?", tradeDay).
		Where("positive_days = total_days").
		Where("trade_count != positive_days").
		Where("reverse_balance != 0").
		Order("balance-total_loss desc").
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

// GetBestForwardCondByTradeDay GetBestForwardCondByTradeDay
func GetBestForwardCondByTradeDay(tradeDay time.Time, db *gorm.DB) (cond simulationcond.AnalyzeCondition, err error) {
	var beforeSort []Result
	err = db.Preload("Cond").Where("trade_day = ?", tradeDay).Where("is_best_forward = true").Find(&beforeSort).Error
	if err != nil {
		return cond, err
	}
	if len(beforeSort) == 0 {
		return cond, err
	}
	return beforeSort[0].Cond, err
}

// GetBestReverseCondByTradeDay GetBestReverseCondByTradeDay
func GetBestReverseCondByTradeDay(tradeDay time.Time, db *gorm.DB) (cond simulationcond.AnalyzeCondition, err error) {
	var beforeSort []Result
	err = db.Preload("Cond").Where("trade_day = ?", tradeDay).Where("is_best_reverse = true").Find(&beforeSort).Error
	if err != nil {
		return cond, err
	}
	if len(beforeSort) == 0 {
		return cond, err
	}
	return beforeSort[0].Cond, err
}

// GetBestCondIDArr GetBestCondIDArr
func GetBestCondIDArr(db *gorm.DB) (idArr []int64, err error) {
	var tmp []Result
	err = db.Where("is_best_forward = true OR is_best_reverse = true").Find(&tmp).Error
	if err != nil {
		return idArr, err
	}
	for _, v := range tmp {
		idArr = append(idArr, v.CondID)
	}
	return idArr, err
}

// GetAllBestForwardResultAndCond GetAllBestForwardResultAndCond
func GetAllBestForwardResultAndCond(db *gorm.DB) (data []Result, err error) {
	err = db.Preload("Cond").Where("is_best_forward = true").Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, err
}

// GetAllBestReverseResultAndCond GetAllBestReverseResultAndCond
func GetAllBestReverseResultAndCond(db *gorm.DB) (data []Result, err error) {
	err = db.Preload("Cond").Where("is_best_reverse = true").Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, err
}

// GetResultByCond GetResultByCond
func GetResultByCond(condID int, db *gorm.DB) (data Result, err error) {
	err = db.Model(&Result{}).Where("cond_id = ?", condID).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, err
}
