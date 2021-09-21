// Package sysparm package sysparm
package sysparm

import "gorm.io/gorm"

// UpdateSysparm UpdateSysparm
func UpdateSysparm(key string, value interface{}, db *gorm.DB) (err error) {
	err = db.Transaction(func(tx *gorm.DB) error {
		if err = tx.Model(&Parameters{}).Where("key = ?", key).Update("value", value).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
