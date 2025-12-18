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

func (q *QuizHistoryRepositoryImpl) preloadQuestionHistory(db *gorm.DB) *gorm.DB {
	return db.
		Preload("QuizHistory").  // Mengambil data QuizHistory dari QuestionHistory
		Preload("AnswerHistory") // Perbaikan typo dari AnswerHIstory -> AnswerHistory
}

func (q *QuizHistoryRepositoryImpl) preloadQuizHistory(db *gorm.DB) *gorm.DB {
	return db.Preload("QuestionHistory") // Mengambil list pertanyaan di dalam QuizHistory
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

// FindAllQuestionHistory implements [IQuizHistoryRepository].
func (q *QuizHistoryRepositoryImpl) FindAllQuestionHistory(ctx context.Context, quizHistoryId uuid.UUID) ([]*models.QuestionHistory, error) {
	var questionHistory []*models.QuestionHistory

	query := q.db.WithContext(ctx).
		Model(&models.QuestionHistory{}).
		Where("quiz_history_id = ?", quizHistoryId)

	// Gunakan preload yang spesifik untuk QuestionHistory
	query = q.preloadQuestionHistory(query)

	if err := query.Order("created_at ASC").Find(&questionHistory).Error; err != nil {
		return nil, err
	}

	return questionHistory, nil
}

// FindQuizHistoryById implements [IQuizHistoryRepository].
func (q *QuizHistoryRepositoryImpl) FindQuizHistoryById(ctx context.Context, quizHistoryId uuid.UUID) (*models.QuizHistory, error) {
	var quizHistory models.QuizHistory

	// Gunakan preload yang spesifik untuk QuizHistory
	query := q.preloadQuizHistory(q.db.WithContext(ctx))

	if err := query.First(&quizHistory, "id = ?", quizHistoryId).Error; err != nil {
		return nil, err
	}

	return &quizHistory, nil
}
