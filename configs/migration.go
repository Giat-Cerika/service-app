package configs

import (
	"giat-cerika-service/internal/models"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Role{},
		&models.Class{},
		&models.User{},
		&models.Image{},
		&models.Materials{},
	)
}
