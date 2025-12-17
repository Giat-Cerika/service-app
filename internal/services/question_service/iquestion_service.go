package questionservice

import (
	"context"
	questionrequest "giat-cerika-service/internal/dto/request/question_request"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type IQuestionService interface {
	CreateQuestion(ctx context.Context, req questionrequest.CreateQuestionRequest) error
	FindQuestionById(ctx context.Context, questionId uuid.UUID) (*models.Question, error)
	FindAllQuestions(ctx context.Context, quizId uuid.UUID, page int, limit int, search string) ([]*models.Question, int, error)
	UpdateQuestion(ctx context.Context, questionId uuid.UUID, req questionrequest.UpdateQuestionRequest) error
	DeleteQuestion(ctx context.Context, questionId uuid.UUID) error
}
