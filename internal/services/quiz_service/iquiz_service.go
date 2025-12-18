package quizservice

import (
	"context"
	quizrequest "giat-cerika-service/internal/dto/request/quiz_request"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type IQuizService interface {
	CreateQuiz(ctx context.Context, req quizrequest.CreateQuizRequest) error
	GetQuizById(ctx context.Context, quizId uuid.UUID) (*models.Quiz, error)
	GetAllQuiz(ctx context.Context, page, limit int, search string) ([]*models.Quiz, int, error)
	UpdateQuiz(ctx context.Context, quizId uuid.UUID, req quizrequest.UpdateQuizRequest) error
	DeleteQuiz(ctx context.Context, quizId uuid.UUID) error
	UpdateStatusQuiz(ctx context.Context, quizId uuid.UUID, req quizrequest.UpdateStatusQuizRequest) error
	UpdateQuestionOrderMode(ctx context.Context, quizId uuid.UUID, req quizrequest.UpdateQuestionOrderModeRequest) error

	GetAllQuizAvailable(ctx context.Context, search string) ([]*models.Quiz, error)
}
