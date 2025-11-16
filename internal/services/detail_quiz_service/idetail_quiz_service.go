package detailquizservice

import (
	"context"
	detailquizrequest "giat-cerika-service/internal/dto/request/detail_quiz_request"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type IQuestionnaireService interface {
	CreateQuestionnaire(ctx context.Context, req detailquizrequest.CreateQuestionnaireRequest) error
	GetAllQuestionnaire(ctx context.Context, page, limit int, search string) ([]*models.Questionnaire, int, error)
	GetByIdQuestionnaire(ctx context.Context, detailquizId uuid.UUID) (*models.Questionnaire, error)
	UpdateQuestionnaire(ctx context.Context, detailquizId uuid.UUID, req detailquizrequest.UpdateQuestionnaireRequest) error
	DeleteQuestionnaire(ctx context.Context, detailquizId uuid.UUID) error
}
