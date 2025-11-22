package configs

import (
	"giat-cerika-service/internal/models"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {

	if err := CreateQuestionnaireEnum(db); err != nil {
		return err
	}

	return db.AutoMigrate(
		&models.Role{},
		&models.Class{},
		&models.User{},
		&models.Image{},
		&models.Materials{},
		&models.MaterialImages{},
		&models.Questionnaire{},
		&models.Video{},
	)
}

func CreateQuestionnaireEnum(db *gorm.DB) error {
	return db.Exec(`
		DO $$ BEGIN
			CREATE TYPE questionnaire_type AS ENUM ('Pengetahuan', 'Perilaku');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
	`).Error
}
