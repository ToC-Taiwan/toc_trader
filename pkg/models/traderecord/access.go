// Package traderecord package traderecord
package traderecord

import (
	"time"

	"gorm.io/gorm"
)

// Insert Insert
func Insert(record TradeRecord, db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertOrUpdate InsertOrUpdate
func InsertOrUpdate(records []TradeRecord, db *gorm.DB) (err error) {
	for _, v := range records {
		if v.OrderID == "" {
			continue
		}
		var exist bool
		price := v.Price
		quantity := v.Quantity
		status := v.Status
		orderID := v.OrderID
		exist, err = CheckExistByOrderID(v.OrderID, db)
		if err != nil {
			return err
		}
		if exist {
			err = db.Transaction(func(tx *gorm.DB) error {
				if err = tx.Model(&TradeRecord{}).Where("order_id = ?", orderID).Updates(map[string]interface{}{
					"price":    price,
					"quantity": quantity,
					"status":   status,
				}).Error; err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				return err
			}
		} else {
			err = Insert(v, db)
			if err != nil {
				return err
			}
		}
	}
	return err
}

// UpdateByStockNumAndClose UpdateByStockNumAndClose
func UpdateByStockNumAndClose(records []*TradeRecord, db *gorm.DB) (err error) {
	for _, v := range records {
		if v.OrderID == "" {
			continue
		}
		var inDB TradeRecord
		inDB, err = GetOrderByOrderID(v.OrderID, db)
		if err != nil {
			return err
		}
		if inDB.ID != 0 {
			record := v
			err = db.Transaction(func(tx *gorm.DB) error {
				if err = tx.Model(&inDB).Updates(record).Error; err != nil {
					return err
				}
				return nil
			})
			continue
		} else {
			err = Insert(*v, db)
		}
	}
	return err
}

// GetCntByOrderID GetCntByOrderID
func GetCntByOrderID(orderID string, db *gorm.DB) (cnt int64, err error) {
	result := db.Model(&TradeRecord{}).Where("order_id = ?", orderID).Count(&cnt)
	if result.Error != nil {
		return cnt, err
	}
	return cnt, err
}

// GetFilledBuyOrder GetFilledBuyOrder
func GetFilledBuyOrder(db *gorm.DB) (order []TradeRecord, err error) {
	result := db.Model(&TradeRecord{}).Where("action = 1 AND status = 6").Find(&order)
	if result.Error != nil {
		return order, err
	}
	return order, err
}

// GetFilledSellOrder GetFilledSellOrder
func GetFilledSellOrder(db *gorm.DB) (order []TradeRecord, err error) {
	result := db.Model(&TradeRecord{}).Where("action = 2 AND status = 6").Find(&order)
	if result.Error != nil {
		return order, err
	}
	return order, err
}

// GetUnFilledOrder GetUnFilledOrder
func GetUnFilledOrder(db *gorm.DB) (order []TradeRecord, err error) {
	result := db.Model(&TradeRecord{}).Where("status = 3").Find(&order)
	if result.Error != nil {
		return order, err
	}
	return order, err
}

// GetAllOrder GetAllOrder
func GetAllOrder(db *gorm.DB) (orderArr []TradeRecord, err error) {
	result := db.Model(&TradeRecord{}).Find(&orderArr)
	if result.Error != nil {
		return orderArr, err
	}
	return orderArr, err
}

// GetAllorderByDayTime GetAllorderByDayTime
func GetAllorderByDayTime(dayTime time.Time, db *gorm.DB) (orderArr []TradeRecord, err error) {
	var tmp []TradeRecord
	result := db.Model(&TradeRecord{}).Order("order_time asc").Find(&tmp)
	if result.Error != nil {
		return orderArr, err
	}
	for _, v := range tmp {
		if v.OrderTime.YearDay() == dayTime.YearDay() {
			orderArr = append(orderArr, v)
		}
	}
	return orderArr, err
}

// CheckExistByOrderID CheckExistByOrderID
func CheckExistByOrderID(orderID string, db *gorm.DB) (exist bool, err error) {
	var cnt int64
	result := db.Model(&TradeRecord{}).Where("order_id = ?", orderID).Count(&cnt)
	if result.Error != nil {
		return exist, err
	}
	if cnt > 0 {
		return true, err
	}
	return exist, err
}

// GetOrderByOrderID GetOrderByOrderID
func GetOrderByOrderID(orderID string, db *gorm.DB) (record TradeRecord, err error) {
	err = db.Model(&TradeRecord{}).Where("order_id = ?", orderID).Find(&record).Error
	return record, err
}

// CheckIsFilledByOrderID CheckIsFilledByOrderID
func CheckIsFilledByOrderID(orderID string, db *gorm.DB) (filled bool, err error) {
	var cnt int64
	err = db.Model(&TradeRecord{}).Where("order_id = ? AND status = 6", orderID).Count(&cnt).Error
	if cnt > 0 {
		return true, err
	}
	return filled, err
}

// GetOrderIDByRecord GetOrderIDByRecord
func GetOrderIDByRecord(record TradeRecord, db *gorm.DB) (orderID string, err error) {
	var tmp TradeRecord
	err = db.Model(&TradeRecord{}).
		Where("stock_num = ?", record.StockNum).
		Where("action = ?", record.Action).
		Where("price = ?", record.Price).
		Where("quantity = ?", record.Quantity).
		Where("status !=5 AND status !=6").
		Find(&tmp).Error
	return tmp.OrderID, err
}
