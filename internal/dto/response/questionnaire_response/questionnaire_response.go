package questionnaireresponse

import (
	"giat-cerika-service/internal/models"
	"giat-cerika-service/pkg/utils"

	"github.com/google/uuid"
)

type QuestionnaireResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Amount      int       `json:"amount"`
	Code        string    `json:"code"`
	Status      int       `json:"status"`
	Type        string    `json:"type"`
	Duration    string    `json:"duration"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

func ToQuestionnaireResponse(questionnaire models.Questionnaire) QuestionnaireResponse {
	return QuestionnaireResponse{
		ID:          questionnaire.ID,
		Title:       questionnaire.Title,
		Description: questionnaire.Description,
		Amount:      questionnaire.Amount,
		Code:        questionnaire.Code,
		Status:      questionnaire.Status,
		Type:        questionnaire.Type,
		Duration:    questionnaire.Duration,
		CreatedAt:   utils.FormatDate(questionnaire.CreatedAt),
		UpdatedAt:   utils.FormatDate(questionnaire.UpdatedAt),
	}
}
