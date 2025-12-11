package quizrepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type IQuizTypeRepository interface {
	Create(ctx context.Context, data *models.QuizType) error
	FindByName(ctx context.Context, name string) (*models.QuizType, error)
	FindAll(ctx context.Context) ([]*models.QuizType, error)
	FindById(ctx context.Context, quizTypeId uuid.UUID) (*models.QuizType, error)
	Update(ctx context.Context, quizTypeId uuid.UUID, data *models.QuizType) error
	Delete(ctx context.Context, quizTypeId uuid.UUID) error
}
