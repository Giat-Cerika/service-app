package quizservice

import (
	"context"
	quizrequest "giat-cerika-service/internal/dto/request/quiz_request"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type IQuizTypeService interface {
	CreateQt(ctx context.Context, req quizrequest.CreateQuizTypeRequest) error
	GetAllQt(ctx context.Context) ([]*models.QuizType, error)
	GetByIdQt(ctx context.Context, quizTypeId uuid.UUID) (*models.QuizType, error)
	UpdateQt(ctx context.Context, quizTypeId uuid.UUID, req quizrequest.UpdateQuizTypeRequest) error
	DeleteQt(ctx context.Context, quizTypeId uuid.UUID) error
}
