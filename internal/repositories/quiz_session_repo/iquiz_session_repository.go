package quizsessionrepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type IQuizSessionRepository interface {
	AssignCodeQuiz(ctx context.Context, quizId uuid.UUID, code string) (uuid.UUID, error)
	SaveQuizSession(ctx context.Context, data *models.QuizSession) error
}
