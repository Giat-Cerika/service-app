package quizresponse

import (
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type QuizResponse struct {
	ID              uuid.UUID `json:"id"`
	QuizType        string    `json:"quiz_type"`
	Code            string    `json:"code"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	StartDate       string    `json:"start_date"`
	EndDate         string    `json:"end_date"`
	Status          int       `json:"status"`
	AmountQuestions int       `json:"amount_questions"`
	AmountAssigned  int       `json:"amount_assigned"`
	CreatedAt       string    `json:"created_at"`
	UpdatedAt       string    `json:"updated_at"`
}

func ToQuizResponse(quiz models.Quiz) QuizResponse {
	return QuizResponse{
		ID:              quiz.ID,
		QuizType:        quiz.QuizType.Name,
		Code:            quiz.Code,
		Title:           quiz.Title,
		Description:     quiz.Description,
		StartDate:       quiz.StartDate.Format("01-02-2006 15:04:05"),
		EndDate:         quiz.EndDate.Format("01-02-2006 15:04:05"),
		Status:          quiz.Status,
		AmountQuestions: quiz.AmountQuestions,
		AmountAssigned:  quiz.AmountAssigned,
		CreatedAt:       quiz.CreatedAt.Format("01-02-2006 15:04:05"),
		UpdatedAt:       quiz.UpdatedAt.Format("01-02-2006 15:04:05"),
	}
}
