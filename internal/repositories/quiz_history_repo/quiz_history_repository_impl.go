package quizhistoryrepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuizHistoryRepositoryImpl struct {
	db *gorm.DB
}

func NewQuizHistoryRepositoryImpl(db *gorm.DB) IQuizHistoryRepository {
	return &QuizHistoryRepositoryImpl{db: db}
}

// FindHistoryByUserID implements [IQuizHistoryRepository].
func (q *QuizHistoryRepositoryImpl) FindHistoryByUserID(ctx context.Context, userId uuid.UUID, search string) ([]*models.QuizHistory, error) {
	var quizHistories []*models.QuizHistory

	query := q.db.WithContext(ctx).Model(&models.QuizHistory{}).Where("user_id = ?", userId)

	if search != "" {
		query = query.Where("title ILIKE ?", "%"+search+"%")
	}

	if err := query.Find(&quizHistories).Error; err != nil {
		return nil, err
	}

	return quizHistories, nil
}
