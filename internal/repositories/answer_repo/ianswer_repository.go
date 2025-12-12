package answerrepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type IAnswerRepository interface {
	CreateAnswer(ctx context.Context, data *models.Answer) error
	DeleteByQuestionID(ctx context.Context, questionID uuid.UUID) error
}
