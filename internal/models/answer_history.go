package models

import (
	"time"

	"github.com/google/uuid"
)

type AnswerHistory struct {
	ID                uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	QuestionHistoryID uuid.UUID       `gorm:"type:uuid;index" json:"question_history_id"`
	QuestionHistory   QuestionHistory `gorm:"foreignKey:QuestionHistoryID;constraint:OnDelete:CASCADE;"`
	AnswerID          uuid.UUID       `gorm:"type:uuid;index" json:"answer_id"`

	AnswerText  string `gorm:"type:text" json:"answer_text"`
	ScoreValue  int    `gorm:"type:int" json:"score_value"`
	ScoreEarned int    `gorm:"type:int" json:"score_earned"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
