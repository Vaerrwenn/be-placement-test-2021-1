package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	SavingID uint
	Type     string `gorm:"size:11"` // DEPOSIT or WITHDRAWAL
	Value    int64
}
