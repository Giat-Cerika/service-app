package questionrepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type IQuestionRepository interface {
	CreateQuestion(ctx context.Context, data *models.Question) error
	FindQuestionById(ctx context.Context, questionId uuid.UUID) (*models.Question, error)
	FindAllQuestions(ctx context.Context, quizId uuid.UUID, limit, offset int, search string) ([]*models.Question, int, error)
	UpdateQuestion(ctx context.Context, questionId uuid.UUID, data *models.Question) error
	DeleteQuestion(ctx context.Context, questionId uuid.UUID) error

	UpdateImageQuestion(ctx context.Context, questionId uuid.UUID, image string) error
}
