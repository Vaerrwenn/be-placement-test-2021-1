package models

import (
	"b-pay/config/database"

	"gorm.io/gorm"
)

// Transaction for each Saving account.
//
// Only has 2 types, DEPOSIT or WITHDRAWAL
type Transaction struct {
	gorm.Model
	SavingID    uint   `gorm:"not null"`
	Type        string `gorm:"size:11;not null;"` // DEPOSIT or WITHDRAWAL
	Value       int64  `gorm:"not null"`
	Description string `gorm:"size:200"`
}

// Store creates a Transaction record to Database.
func (t *Transaction) Store() error {
	err := database.DB.Create(&t).Error
	return err
}
