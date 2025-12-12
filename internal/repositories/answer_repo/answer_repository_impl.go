package answerrepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnswerRepositoryImpl struct {
	db *gorm.DB
}

func NewAnswerRepositoryImpl(db *gorm.DB) IAnswerRepository {
	return &AnswerRepositoryImpl{db: db}
}

// CreateAnswer implements IAnswerRepository.
func (a *AnswerRepositoryImpl) CreateAnswer(ctx context.Context, data *models.Answer) error {
	return a.db.WithContext(ctx).Create(data).Error
}

// DeleteByQuestionID implements IAnswerRepository.
func (a *AnswerRepositoryImpl) DeleteByQuestionID(ctx context.Context, questionID uuid.UUID) error {
	return a.db.WithContext(ctx).Where("question_id = ?", questionID).Delete(&models.Answer{}).Error
}
