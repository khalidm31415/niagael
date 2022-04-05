package entity

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&Order{}); err != nil {
			return err
		}
		if err := tx.AutoMigrate(&OrderItem{}); err != nil {
			return err
		}
		return nil
	})
	return err
}
