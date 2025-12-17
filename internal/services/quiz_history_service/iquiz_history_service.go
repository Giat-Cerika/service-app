package quizhistoryservice

import (
	"context"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type IQuizHistoryService interface {
	GetHistoryQuizStudent(ctx context.Context, userId uuid.UUID, search string) ([]*models.QuizHistory, error)
}
