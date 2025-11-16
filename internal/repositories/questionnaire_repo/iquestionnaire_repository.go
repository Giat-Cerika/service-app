package questionnairerepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type IQuestionnaireRepository interface {
	Create(ctx context.Context, data *models.Questionnaire) error
	FindByTitle(ctx context.Context, NameQuestionnaire string) (*models.Questionnaire, error)
	FindAll(ctx context.Context, limit, offset int, search string) ([]*models.Questionnaire, int, error)
	FindById(ctx context.Context, questionnaireId uuid.UUID) (*models.Questionnaire, error)
	Update(ctx context.Context, questionnaireId uuid.UUID, data *models.Questionnaire) error
	Delete(ctx context.Context, questionnaireId uuid.UUID) error
}
