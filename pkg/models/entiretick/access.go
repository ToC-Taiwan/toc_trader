// Package entiretick package entiretick
package entiretick

import (
	"time"

	"gitlab.tocraw.com/root/toc_trader/global"
	"gorm.io/gorm"
)

// Insert Insert
func Insert(tick EntireTick, db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&tick).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertMultiRecord InsertMultiRecord
func InsertMultiRecord(tickArr []*EntireTick, db *gorm.DB) error {
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

// GetCntByStockAndDate GetCntByStockAndDate
func GetCntByStockAndDate(stockNum, date string, db *gorm.DB) (cnt int64, err error) {
	utcTimeSec, err := time.Parse(global.LongTimeLayout, date+" 00:00:00")
	if err != nil {
		return cnt, err
	}
	startTime := utcTimeSec.UnixNano()
	endTime := utcTimeSec.AddDate(0, 0, 1).UnixNano()
	result := db.Model(&EntireTick{}).Where("stock_num = ? AND timestamp >= ? AND timestamp < ?", stockNum, startTime, endTime).Count(&cnt)
	return cnt, result.Error
}

// GetLastCloseByDate GetLastCloseByDate
func GetLastCloseByDate(stockNum, date string, db *gorm.DB) (close float64, err error) {
	utcTimeSec, err := time.Parse(global.LongTimeLayout, date+" 00:00:00")
	if err != nil {
		return close, err
	}
	startTime := utcTimeSec.UnixNano()
	endTime := utcTimeSec.AddDate(0, 0, 1).UnixNano()

	var data EntireTick
	result := db.Model(&EntireTick{}).Where("stock_num = ? AND timestamp >= ? AND timestamp < ?", stockNum, startTime, endTime).
		Order("timestamp desc").Limit(1).Find(&data)
	return data.Close, result.Error
}

// GetByStockAndTimeStamp GetByStockAndTimeStamp
func GetByStockAndTimeStamp(stockNum string, timestamp int64, db *gorm.DB) (records []EntireTick, err error) {
	result := db.Model(&EntireTick{}).Where("stock_num = ? AND timestamp >= ?", stockNum, timestamp).Order("timestamp asc").Find(&records)
	return records, result.Error
}

// GetAllEntiretickByStock GetAllEntiretickByStock
func GetAllEntiretickByStock(stockNum string, db *gorm.DB) (records []*EntireTick, err error) {
	result := db.Model(&EntireTick{}).Where("stock_num = ?", stockNum).Order("timestamp asc").Find(&records)
	return records, result.Error
}

// GetAllEntiretickByStockByDate GetAllEntiretickByStockByDate
func GetAllEntiretickByStockByDate(stockNum, date string, db *gorm.DB) (records []*EntireTick, err error) {
	utcTimeSec, err := time.Parse(global.LongTimeLayout, date+" 00:00:00")
	if err != nil {
		return records, err
	}
	startTime := utcTimeSec.UnixNano()
	endTime := utcTimeSec.AddDate(0, 0, 1).UnixNano()

	result := db.Model(&EntireTick{}).Where("stock_num = ? AND timestamp >= ? AND timestamp < ?", stockNum, startTime, endTime).Order("timestamp asc").Find(&records)
	return records, result.Error
}

// DeleteAll DeleteAll
func DeleteAll(db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Not("id = 0").Unscoped().Delete(&EntireTick{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteByStockNumAndDate DeleteByStockNumAndDate
func DeleteByStockNumAndDate(stockNum, date string, db *gorm.DB) error {
	utcTimeSec, err := time.Parse(global.LongTimeLayout, date+" 00:00:00")
	if err != nil {
		return err
	}
	startTime := utcTimeSec.UnixNano()
	endTime := utcTimeSec.AddDate(0, 0, 1).UnixNano()
	err = db.Transaction(func(tx *gorm.DB) error {
		if err = tx.Where("stock_num = ? AND timestamp >= ? AND timestamp < ?", stockNum, startTime, endTime).Unscoped().Delete(&EntireTick{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
