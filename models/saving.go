package models

import (
	"b-pay/config/database"

	"gorm.io/gorm"
)

// Saving defines every saving's data.
type Saving struct {
	gorm.Model
	UserID       uint
	Name         string `gorm:"size:100"`
	Balance      int64
	PIN          []byte `gorm:"size:6"`
	Transactions []Transaction
}

// SavingIndex is a struct for GetSavingsByUserID return value.
type SavingIndex struct {
	ID      int
	Name    string
	Balance string
}

// Store stores Saving data to DB.
func (s *Saving) Store() error {
	err := database.DB.Create(&s).Error
	return err
}

// GetSavingsByUserID get/fetch multiple Saving data with corresponded userID.
func (s *Saving) GetSavingsByUserID(userID string) (*[]SavingIndex, error) {
	var results []SavingIndex
	query := database.DB.Model(&Saving{}).
		Select("id, name, balance").
		Where("user_id = ?", userID).
		Scan(&results)

	if query.Error != nil {
		return nil, query.Error
	}
	return &results, nil
}

// GetPINBySavingID gets/fetches a Saving PIN by searching Saving ID.
func (s *Saving) GetPINBySavingID(savingID string) string {
	var result string
	err := database.DB.Model(&Saving{}).
		Select("pin").
		Where("id = ?", savingID).
		First(&result).
		Error

	if err != nil {
		return ""
	}

	return result
}

// GetSavingByID gets/fetches Saving data by searching the ID.
func (s *Saving) GetSavingByID(id string) *Saving {
	var result Saving
	err := database.DB.Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil
	}
	return &result
}

// Update updates the "source" data with the inputted data.
func (s *Saving) Update(source *Saving) error {
	err := database.DB.Model(&source).Updates(&s).Error
	return err
}

// Delete deletes a Saving account data.
// Changes balance to 0 first before soft-deleting the data.
func (s *Saving) Delete() error {
	err := database.DB.Model(&s).Update("balance", 0).Error
	if err != nil {
		return err
	}

	err = database.DB.Delete(&s).Error
	return err
}
