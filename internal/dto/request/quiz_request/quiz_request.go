package quizrequest

import (
	"time"

	"github.com/google/uuid"
)

type CreateQuizRequest struct {
	QuizTypeId  uuid.UUID `form:"quiz_type_id" json:"quiz_type_id"`
	Code        string    `form:"code" json:"code"`
	Title       string    `form:"title" json:"title"`
	Description string    `form:"description" json:"description"`
	StartDate   time.Time `form:"start_date" json:"start_date"`
	EndDate     time.Time `form:"end_date" json:"end_date"`
}

type UpdateQuizRequest struct {
	QuizTypeId  uuid.UUID `form:"quiz_type_id" json:"quiz_type_id"`
	Code        string    `form:"code" json:"code"`
	Title       string    `form:"title" json:"title"`
	Description string    `form:"description" json:"description"`
	StartDate   time.Time `form:"start_date" json:"start_date"`
	EndDate     time.Time `form:"end_date" json:"end_date"`
}

type UpdateStatusQuizRequest struct {
	Status int `form:"status" json:"status"`
}

type UpdateQuestionOrderModeRequest struct {
	QuestionOrderMode string `form:"question_order_mode" json:"question_order_mode"`
}
