package models

import (
	"time"

	"github.com/google/uuid"
)

type Answer struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	QuestionID uuid.UUID `gorm:"type:uuid;index" json:"question_id"`
	Question   Question  `gorm:"foreignKey:QuestionID"`
	AnswerText string    `gorm:"type:text" json:"answer_text"`
	ScoreValue int       `gorm:"type:int" json:"score_value"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
