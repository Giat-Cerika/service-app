package models

import (
	"time"

	"github.com/google/uuid"
)

type QuizSessionStatus string

const (
	SessionStatusStarted    QuizSessionStatus = "started"
	SessionStatusInProgress QuizSessionStatus = "in_progress"
	SessionStatusCompleted  QuizSessionStatus = "completed"
)

type QuizSession struct {
	ID          uuid.UUID         `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID      uuid.UUID         `gorm:"type:uuid;index" json:"user_id"`
	User        User              `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	QuizID      uuid.UUID         `gorm:"type:uuid;index" json:"quiz_id"`
	Quiz        Quiz              `gorm:"foreignKey:QuizID;contstarint:OnDelete:CASCADE;"`
	Score       int               `gorm:"type:int" json:"score"`
	MaxScore    int               `gorm:"type:int" json:"max_score"`
	Status      QuizSessionStatus `gorm:"type:varchar(50);default:'started'" json:"status"`
	StartedAt   time.Time         `gorm:"type:timestamp" json:"started_at"`
	CompletedAt time.Time         `gorm:"type:timestamp" json:"completed_at"`
	CreatedAt   time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time         `gorm:"autoUpdateTime" json:"updated_at"`

	Responses []Response `gorm:"constraint:OnDelete:CASCADE;"`
}
