package quizresponse

import (
	"giat-cerika-service/internal/models"
	"giat-cerika-service/pkg/utils"

	"github.com/google/uuid"
)

type QuizTypeResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

func ToQuizTypeResponse(qt models.QuizType) QuizTypeResponse {
	return QuizTypeResponse{
		ID:          qt.ID,
		Name:        qt.Name,
		Description: qt.Description,
		CreatedAt:   utils.FormatDate(qt.CreatedAt),
		UpdatedAt:   utils.FormatDate(qt.UpdatedAt),
	}
}
