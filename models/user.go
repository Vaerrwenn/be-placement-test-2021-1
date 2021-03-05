package models

import (
	"b-pay/config/database"

	"gorm.io/gorm"
)

// User defines every user's data.
type User struct {
	gorm.Model
	Name     string `gorm:"size:100;not null;"`
	Email    string `gorm:"size:300;unique;not null;"`
	Password []byte `gorm:"not null"`
	Savings  []Saving
}

// StoreUser stores User data into Database.
func (u *User) StoreUser() error {
	err := database.DB.Create(&u).Error
	return err
}

// GetUserByEmail searches a User by presented Email.
// Returns the User data.
func (u *User) GetUserByEmail() *User {
	var result = &User{}
	err := database.DB.Where(map[string]interface{}{
		"email": u.Email,
	}).First(&result).Error
	if err != nil {
		return nil
	}
	return result
}
