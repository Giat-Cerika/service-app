package models

import (
	"time"

	"github.com/google/uuid"
)

type QuestionHistory struct {
	ID            uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	QuizHistoryID uuid.UUID   `gorm:"type:uuid;index" json:"quiz_history_id"`
	QuizHistory   QuizHistory `gorm:"foreignKey:QuizHistoryID;constraint:OnDelete:CASCADE;"`
	QuestionID    uuid.UUID   `gorm:"type:uuid;index" json:"question_id"`

	QuestionText  string `gorm:"type:text" json:"question_text"`
	QuestionImage string `gorm:"type:varchar(255); null" json:"question_image"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
