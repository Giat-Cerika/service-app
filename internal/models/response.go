package models

import (
	"time"

	"github.com/google/uuid"
)

type Response struct {
	ID            uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	QuizSessionID uuid.UUID   `gorm:"type:uuid;index" json:"quiz_session_id"`
	QuizSession   QuizSession `gorm:"foreignKey:QuizSessionID;constraint:OnDelete:CASCADE;"`
	QuestionID    uuid.UUID   `gorm:"type:uuid;index" json:"question_id"`
	Question      Question    `gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE;"`
	AnswerID      *uuid.UUID  `gorm:"type:uuid; null" json:"answer_id"`
	ScoreEarned   int         `gorm:"type:int" json:"score_earned"`
	CreatedAt     time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}
