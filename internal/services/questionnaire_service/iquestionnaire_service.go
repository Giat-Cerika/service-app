package questionnaireservice

import (
	"context"
	questionnairerequest "giat-cerika-service/internal/dto/request/questionnaire_request"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type IQuestionnaireService interface {
	CreateQuestionnaire(ctx context.Context, req questionnairerequest.CreateQuestionnaireRequest) error
	GetAllQuestionnaire(ctx context.Context, page, limit int, search string) ([]*models.Questionnaire, int, error)
	GetByIdQuestionnaire(ctx context.Context, questionnaireId uuid.UUID) (*models.Questionnaire, error)
	UpdateQuestionnaire(ctx context.Context, questionnaireId uuid.UUID, req questionnairerequest.UpdateQuestionnaireRequest) error
	DeleteQuestionnaire(ctx context.Context, questionnaireId uuid.UUID) error
}
