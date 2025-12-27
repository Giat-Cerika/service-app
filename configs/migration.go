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
		&models.Video{},
		&models.ToootBrushLog{},
		&models.QuizType{},
		&models.Quiz{},
		&models.Question{},
		&models.Answer{},
		&models.QuizSession{},
		&models.Response{},
		&models.QuizHistory{},
		&models.QuestionHistory{},
		&models.AnswerHistory{},
		&models.StimulatedSaliva{},
		&models.RestingSaliva{},
		&models.SalivaOption{},
		&models.PlaqueOption{},
		&models.DietDetail{},
		&models.CariesRisk{},
		&models.ConfidenceDetail{},
		&models.Prediction{},
		&models.PredictHistory{},
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
