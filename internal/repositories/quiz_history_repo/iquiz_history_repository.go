package quizhistoryrepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type IQuizHistoryRepository interface {
	FindHistoryByUserID(ctx context.Context, userId uuid.UUID, search string) ([]*models.QuizHistory, error)
}
