// Package stock package stock
package stock

import (
	"time"

	"gitlab.tocraw.com/root/toc_trader/pkg/models/sysparm"
	"gorm.io/gorm"
)

// Insert Insert
func Insert(stock Stock, db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&stock).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertMultiRecord InsertMultiRecord
func InsertMultiRecord(stockArr []Stock, db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		batch := len(stockArr)
		if batch >= 2000 {
			batch = 2000
		}
		if err := tx.CreateInBatches(&stockArr, batch).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// GetTargetByLowHighVolume GetTargetByLowHighVolume
func GetTargetByLowHighVolume(low, high, volume int64, db *gorm.DB) (data []Stock, err error) {
	result := db.Model(&Stock{}).Where("day_trade = ? AND last_close > ? AND last_close < ? AND last_volume > ?", true, low, high, volume).
		Order("stock_num").Find(&data)
	return data, result.Error
}

// GetStocksFromNumArr GetStocksFromNumArr
func GetStocksFromNumArr(stockNumArr []string, db *gorm.DB) (data []Stock, err error) {
	for _, v := range stockNumArr {
		var tmpStock Stock
		err = db.Model(&Stock{}).Where("stock_num = ?", v).Find(&tmpStock).Error
		data = append(data, tmpStock)
	}
	return data, err
}

// GetTargetByMultiLowHighVolume GetTargetByMultiLowHighVolume
func GetTargetByMultiLowHighVolume(conditionArr []sysparm.TargetCondArr, db *gorm.DB) (data []Stock, err error) {
	for _, cond := range conditionArr {
		var tmp []Stock
		result := db.Model(&Stock{}).Where("day_trade = ? AND last_close > ? AND last_close < ? AND last_volume > ?", true, cond.LimitPriceLow, cond.LimitPriceHigh, cond.LimitVolume).Order("stock_num").Find(&tmp)
		if result.Error != nil {
			return data, err
		}
		data = append(data, tmp...)
	}
	return data, err
}

// GetTotalCnt GetTotalCnt
func GetTotalCnt(db *gorm.DB) (cnt int64, err error) {
	result := db.Model(&Stock{}).Count(&cnt)
	if result.Error != nil {
		return cnt, err
	}
	return cnt, err
}

// GetLastUpdatedTime GetLastUpdatedTime
func GetLastUpdatedTime(db *gorm.DB) (updatedTime time.Time, err error) {
	var data Stock
	result := db.Model(&Stock{}).Where("last_volume > 0").Not("updated_at IS NULL").Order("updated_at").Limit(1).Find(&data)
	if result.Error != nil {
		return updatedTime, err
	}
	return data.UpdatedAt, err
}

// GetAllStockNum GetAllStockNum
func GetAllStockNum(db *gorm.DB) (allStockNum []string, err error) {
	var data []Stock
	var resultArr []string
	result := db.Model(&Stock{}).Not("stock_name = ?", "").Order("stock_num").Find(&data)
	if result.Error != nil {
		return allStockNum, err
	}
	for _, v := range data {
		resultArr = append(resultArr, v.StockNum)
	}
	return resultArr, err
}

// GetAllRows GetAllRows
func GetAllRows(db *gorm.DB) (allStock []Stock, err error) {
	result := db.Model(&Stock{}).Not("stock_name = ?", "").Order("stock_num").Find(&allStock)
	if result.Error != nil {
		return allStock, err
	}
	return allStock, err
}

// CheckExistByStockNum CheckExistByStockNum
func CheckExistByStockNum(stockNum string, db *gorm.DB) (exist bool, err error) {
	var cnt int64
	result := db.Model(&Stock{}).Where("stock_num = ?", stockNum).Count(&cnt)
	if result.Error != nil {
		return exist, err
	}
	if cnt != 0 {
		return true, err
	}
	return false, err
}
