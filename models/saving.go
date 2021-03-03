package models

import "gorm.io/gorm"

// Saving defines every saving's data.
type Saving struct {
	gorm.Model
	UserID       uint
	Balance      int64
	Password     []byte
	Transactions []Transaction
}
