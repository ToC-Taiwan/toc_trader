// Package holiday package holiday
package holiday

import "gorm.io/gorm"

// Insert Insert
func Insert(holiday Holiday, db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&holiday).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertMultiRecord InsertMultiRecord
func InsertMultiRecord(holidayArr []Holiday, db *gorm.DB) error {
	batch := len(holidayArr)
	if batch >= 2000 {
		batch = 2000
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(&holidayArr, batch).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// GetAllHoliday GetAllHoliday
func GetAllHoliday(db *gorm.DB) (allHoliday []Holiday, err error) {
	result := db.Model(&Holiday{}).Find(&allHoliday)
	if result.Error != nil {
		return allHoliday, err
	}
	return allHoliday, err
}

// CheckIsHolidayByTimeStamp CheckIsHolidayByTimeStamp
func CheckIsHolidayByTimeStamp(timeStamp int64, db *gorm.DB) (exist bool, err error) {
	var cnt int64
	result := db.Model(&Holiday{}).Where("timestamp = ?", timeStamp).Count(&cnt)
	if result.Error != nil {
		return exist, err
	}
	if cnt != 0 {
		return true, err
	}
	return exist, err
}
