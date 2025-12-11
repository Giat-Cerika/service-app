package quizrepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuizTypeRepositoryImpl struct {
	db *gorm.DB
}

func NewQuizTypeRepositoryImpl(db *gorm.DB) IQuizTypeRepository {
	return &QuizTypeRepositoryImpl{db: db}
}

// Create implements IQuizTypeRepository.
func (q *QuizTypeRepositoryImpl) Create(ctx context.Context, data *models.QuizType) error {
	return q.db.WithContext(ctx).Create(data).Error
}

// FindByName implements IQuizTypeRepository.
func (q *QuizTypeRepositoryImpl) FindByName(ctx context.Context, name string) (*models.QuizType, error) {
	var qt models.QuizType
	if err := q.db.WithContext(ctx).First(&qt, "name = ?", name).Error; err != nil {
		return nil, err
	}

	return &qt, nil
}

// FindAll implements IQuizTypeRepository.
func (q *QuizTypeRepositoryImpl) FindAll(ctx context.Context) ([]*models.QuizType, error) {
	var qts []*models.QuizType
	if err := q.db.WithContext(ctx).Order("created_at DESC").Find(&qts).Error; err != nil {
		return nil, err
	}

	return qts, nil
}

// FindById implements IQuizTypeRepository.
func (q *QuizTypeRepositoryImpl) FindById(ctx context.Context, quizTypeId uuid.UUID) (*models.QuizType, error) {
	var qt models.QuizType
	if err := q.db.WithContext(ctx).First(&qt, "id = ?", quizTypeId).Error; err != nil {
		return nil, err
	}
	return &qt, nil
}

// Update implements IQuizTypeRepository.
func (q *QuizTypeRepositoryImpl) Update(ctx context.Context, quizTypeId uuid.UUID, data *models.QuizType) error {
	return q.db.WithContext(ctx).Save(data).Error
}

// Delete implements IQuizTypeRepository.
func (q *QuizTypeRepositoryImpl) Delete(ctx context.Context, quizTypeId uuid.UUID) error {
	return q.db.WithContext(ctx).Delete(&models.QuizType{}, "id = ?", quizTypeId).Error
}
