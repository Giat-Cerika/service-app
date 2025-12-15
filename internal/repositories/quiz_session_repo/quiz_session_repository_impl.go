package quizsessionrepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuizSessionRepositoryImpl struct {
	db *gorm.DB
}

func NewQuizSessionRepositoryImpl(db *gorm.DB) IQuizSessionRepository {
	return &QuizSessionRepositoryImpl{db: db}
}

// AssignCodeQuiz implements [IQuizSessionRepository].
func (q *QuizSessionRepositoryImpl) AssignCodeQuiz(ctx context.Context, quizId uuid.UUID, code string) (uuid.UUID, error) {
	var quiz models.Quiz

	err := q.db.WithContext(ctx).Where("id = ? AND code = ?", quizId, code).First(&quiz).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return uuid.Nil, nil
		}
		return uuid.Nil, nil
	}

	return quiz.ID, nil
}

// SaveQuizSession implements [IQuizSessionRepository].
func (q *QuizSessionRepositoryImpl) SaveQuizSession(ctx context.Context, data *models.QuizSession) error {
	return q.db.WithContext(ctx).Create(data).Error
}
