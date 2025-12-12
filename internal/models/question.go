package models

import (
	"time"

	"github.com/google/uuid"
)

type Question struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	QuizID        uuid.UUID `gorm:"type:uuid;index" json:"quiz_id"`
	Quiz          Quiz      `gorm:"foreignKey:QuizID"`
	QuestionText  string    `gorm:"type:text" json:"question_text"`
	QuestionImage string    `gorm:"type:varchar(255); null" json:"question_image"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Answers []Answer `gorm:"constraint:OnDelete:CASCADE;"`
}
