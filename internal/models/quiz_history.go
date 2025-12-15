package models

import (
	"time"

	"github.com/google/uuid"
)

type QuizHistory struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	QuizID          uuid.UUID `gorm:"type:uuid;index" json:"quiz_id"`
	QuizSessionID   uuid.UUID `gorm:"type:uuid;index" json:"quiz_session_id"`
	Code            string    `gorm:"type:varchar(100);index" json:"code"`
	Title           string    `gorm:"type:varchar(255)" json:"title"`
	Description     string    `gorm:"type:text" json:"description"`
	StartDate       time.Time `gorm:"type:timestamp" json:"start_date"`
	EndDate         time.Time `gorm:"type:timestamp" json:"end_date"`
	AmountQuestions int       `gorm:"type:int" json:"amount_questions"`
	AmountAssigned  int       `gorm:"type:int" json:"amount_assigned"`

	UserID      uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	Score       int       `gorm:"type:int" json:"score"`
	MaxScore    int       `gorm:"type:int" json:"max_score"`
	Percentage  float64   `gorm:"type:float" json:"percentage"`
	Status      int       `gorm:"type:int" json:"status"`
	StartedAt   time.Time `gorm:"type:timestamp" json:"started_at"`
	CompletedAt time.Time `gorm:"type:timestamp" json:"completed_at"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
