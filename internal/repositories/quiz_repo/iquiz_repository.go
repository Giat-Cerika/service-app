package quizrepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type IQuizRepository interface {
	Create(ctx context.Context, data *models.Quiz) error
	FindById(ctx context.Context, quizId uuid.UUID) (*models.Quiz, error)
	FindAll(ctx context.Context, limit, offset int, search string) ([]*models.Quiz, int, error)
	Update(ctx context.Context, quizId uuid.UUID, data *models.Quiz) error
	Delete(ctx context.Context, quizId uuid.UUID) error
	UpdateStatus(ctx context.Context, quizId uuid.UUID, status int) error
	IncreamentAmountQuestion(ctx context.Context, quizId uuid.UUID) error
	DecreaseAmountQuestion(ctx context.Context, quizId uuid.UUID) error
}
