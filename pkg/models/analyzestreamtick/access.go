// Package analyzestreamtick analyzestreamtick
package analyzestreamtick

import "gorm.io/gorm"

// Insert Insert
func Insert(tick AnalyzeStreamTick, db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&tick).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertMultiRecord InsertMultiRecord
func InsertMultiRecord(tickArr []*AnalyzeStreamTick, db *gorm.DB) error {
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

// GetLastNRows GetLastNRows
func GetLastNRows(stockNum, n int64, db *gorm.DB) (data []AnalyzeStreamTick, err error) {
	result := db.Model(&AnalyzeStreamTick{}).Where("stock_num = ?", stockNum).Order("timestamp desc").Limit(int(n)).Find(&data)
	if result.Error != nil {
		return data, err
	}
	return data, err
}
