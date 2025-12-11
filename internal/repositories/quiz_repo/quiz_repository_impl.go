package quizrepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuizRepositoryImpl struct {
	db *gorm.DB
}

func NewQuizRepositoryImpl(db *gorm.DB) IQuizRepository {
	return &QuizRepositoryImpl{db: db}
}

// Create implements IQuizRepository.
func (q *QuizRepositoryImpl) Create(ctx context.Context, data *models.Quiz) error {
	return q.db.WithContext(ctx).Create(data).Error
}

// FindAll implements IQuizRepository.
func (q *QuizRepositoryImpl) FindAll(ctx context.Context, limit, offset int, search string) ([]*models.Quiz, int, error) {
	var (
		quizzes []*models.Quiz
		count   int64
	)

	query := q.db.WithContext(ctx).Model(&models.Quiz{})
	if search != "" {
		query = query.Where("title ILIKE ?", "%"+search+"%")
	}
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("QuizType").Limit(limit).Offset(offset).Order("created_at DESC").Find(&quizzes).Error; err != nil {
		return nil, 0, err
	}

	return quizzes, int(count), nil
}

// FindById implements IQuizRepository.
func (q *QuizRepositoryImpl) FindById(ctx context.Context, quizId uuid.UUID) (*models.Quiz, error) {
	var quiz models.Quiz
	if err := q.db.WithContext(ctx).Preload("QuizType").First(&quiz, "id = ?", quizId).Error; err != nil {
		return nil, err
	}
	return &quiz, nil
}

// Update implements IQuizRepository.
func (q *QuizRepositoryImpl) Update(ctx context.Context, quizId uuid.UUID, data *models.Quiz) error {
	return q.db.WithContext(ctx).Save(data).Error
}

// Delete implements IQuizRepository.
func (q *QuizRepositoryImpl) Delete(ctx context.Context, quizId uuid.UUID) error {
	return q.db.WithContext(ctx).Delete(&models.Quiz{}, "id = ?", quizId).Error
}

// UpdateStatus implements IQuizRepository.
func (q *QuizRepositoryImpl) UpdateStatus(ctx context.Context, quizId uuid.UUID, status int) error {
	return q.db.WithContext(ctx).Model(&models.Quiz{}).Where("id = ?", quizId).Update("status", status).Error
}
