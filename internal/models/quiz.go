package models

import (
	"time"

	"github.com/google/uuid"
)

type QuestionOrderMode string

const (
	QuestionOrderSequential QuestionOrderMode = "sequential"
	QuestionOrderRandom     QuestionOrderMode = "random"
)

type Quiz struct {
	ID                uuid.UUID         `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	QuizTypeID        uuid.UUID         `gorm:"type:uuid"`
	QuizType          QuizType          `gorm:"foreignKey:QuizTypeID"`
	Code              string            `gorm:"type:varchar(255);index" json:"code"`
	Title             string            `gorm:"type:varchar(255)" json:"title"`
	Description       string            `gorm:"type:text" json:"description"`
	StartDate         time.Time         `gorm:"type:timestamp" json:"start_date"`
	EndDate           time.Time         `gorm:"type:timestamp" json:"end_date"`
	Status            int               `gorm:"type:int" json:"status"`
	AmountQuestions   int               `gorm:"type:int" json:"amount_questions"`
	AmountAssigned    int               `gorm:"type:int" json:"amount_assigned"`
	QuestionOrderMode QuestionOrderMode `gorm:"type:varchar(50);default:'sequential'" json:"question_order_mode"`
	CreatedAt         time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time         `gorm:"autoUpdateTime" json:"updated_at"`

	Questions []Question `gorm:"constraint:OnDelete:CASCADE;"`
}
