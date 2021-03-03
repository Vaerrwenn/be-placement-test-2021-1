package migration

import (
	"b-pay/models"

	"gorm.io/gorm"
)

// AutoMigrate uses GORM's AutoMigrate to migrate Models.
// Input models manually
func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&models.User{},
		&models.Saving{},
		&models.Transaction{},
	)
}
