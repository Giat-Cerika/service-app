package models

import (
	"time"

	"github.com/google/uuid"
)

type PredictHistory struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PredictionID uuid.UUID  `gorm:"type:uuid" json:"prediction_id"`
	Prediction   Prediction `gorm:"foreignKey:PredictionID;constraint:OnDelete:CASCADE;"`
	UserID       uuid.UUID  `gorm:"type:uuid" json:"user_id"`
	User         User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Suggestion   string     `gorm:"type:text"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}
