package quizhistoryrepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type IQuizHistoryRepository interface {
	FindHistoryByUserID(ctx context.Context, userId uuid.UUID, search string) ([]*models.QuizHistory, error)
	FindAllQuestionHistory(ctx context.Context, quizHistoryId uuid.UUID) ([]*models.QuestionHistory, error)
	FindQuizHistoryById(ctx context.Context, quizHistoryId uuid.UUID) (*models.QuizHistory, error)
	FindHistoryByQuizID(ctx context.Context) ([]*models.QuizHistory, error)
}
