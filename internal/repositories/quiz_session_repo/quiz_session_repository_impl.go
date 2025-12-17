package quizsessionrepo

import (
	"context"
	"giat-cerika-service/internal/models"
	"time"

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

// CreateStartedAt implements [IQuizSessionRepository].
func (q *QuizSessionRepositoryImpl) CreateStartedAt(ctx context.Context, quizSessionId uuid.UUID) error {
	now := time.Now()
	return q.db.WithContext(ctx).
		Model(&models.QuizSession{}).
		Where("id = ?", quizSessionId).
		Updates(map[string]interface{}{
			"started_at": now,
			"status":     models.SessionStatusInProgress,
		}).Error
}

// FindById implements [IQuizSessionRepository].
func (q *QuizSessionRepositoryImpl) FindById(ctx context.Context, userId uuid.UUID, quizSessionId uuid.UUID) (*models.QuizSession, error) {
	var qs models.QuizSession
	if err := q.db.WithContext(ctx).First(&qs, "id = ? and user_id = ?", quizSessionId, userId).Error; err != nil {
		return nil, err
	}

	return &qs, nil
}

// FindByUserAndQuiz implements [IQuizSessionRepository].
func (q *QuizSessionRepositoryImpl) FindByUserAndQuiz(ctx context.Context, userId uuid.UUID, quizId uuid.UUID) (*models.QuizSession, error) {
	var qs models.QuizSession
	err := q.db.WithContext(ctx).
		Where("user_id = ? AND quiz_id = ? AND status IN ?", userId, quizId, []models.QuizSessionStatus{
			models.SessionStatusStarted,
			models.SessionStatusInProgress,
		}).
		Order("created_at DESC").
		First(&qs).Error

	if err != nil {
		return nil, err
	}

	return &qs, nil
}

// SubmitQuiz implements [IQuizSessionRepository].
func (q *QuizSessionRepositoryImpl) BulkSaveResponses(ctx context.Context, quizSessionId uuid.UUID, data []*models.Response) error {
	return q.db.WithContext(ctx).Create(data).Error
}

// CompleteQuizSession implements [IQuizSessionRepository].
func (q *QuizSessionRepositoryImpl) CompleteQuizSession(ctx context.Context, quizSessionId uuid.UUID, score int, maxScore int, completedAt *time.Time) error {
	return q.db.WithContext(ctx).
		Model(&models.QuizSession{}).
		Where("id = ?", quizSessionId).
		Updates(map[string]interface{}{
			"score":        score,
			"max_score":    maxScore,
			"completed_at": completedAt,
			"status":       models.SessionStatusCompleted,
		}).Error
}

// SaveQuizHistoryTransaction implements [IQuizSessionRepository].
func (q *QuizSessionRepositoryImpl) SaveQuizHistoryTransaction(ctx context.Context, quizHistory *models.QuizHistory, questionHistories []models.QuestionHistory, answerHistories []models.AnswerHistory) error {
	return q.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 1. Simpan QuizHistory
		if err := tx.Create(quizHistory).Error; err != nil {
			return err // Akan memicu rollback
		}

		// 2. Simpan QuestionHistories (Bulk Insert)
		if len(questionHistories) > 0 {
			// Karena QuestionHistory adalah relasi Many-to-One ke QuizHistory,
			// dan ID QuizHistory sudah ada, kita bisa langsung simpan.
			if err := tx.Create(&questionHistories).Error; err != nil {
				return err // Akan memicu rollback
			}
		} else {
			// Seharusnya tidak terjadi, tapi jika tidak ada pertanyaan
			return gorm.ErrInvalidData
		}

		// 3. Simpan AnswerHistories (Bulk Insert)
		if len(answerHistories) > 0 {
			// Karena AnswerHistory adalah relasi Many-to-One ke QuestionHistory,
			// dan ID QuestionHistory sudah ada, kita bisa langsung simpan.
			// Bulk insert AnswerHistory mungkin memakan waktu lama jika kuis besar.
			// Anda mungkin ingin menggunakan CreateInBatches(data, N) untuk performa.
			if err := tx.Create(&answerHistories).Error; err != nil {
				return err // Akan memicu rollback
			}
		}

		return nil
	})
}

// FindQuizWithOrderedQuestions implements [IQuizSessionRepository].
func (q *QuizSessionRepositoryImpl) FindQuizWithOrderedQuestions(ctx context.Context, quizId uuid.UUID, orderMode string) (*models.Quiz, error) {
	var quiz models.Quiz
	var orderQuery string

	if orderMode == "random" {
		orderQuery = "RANDOM()"
	} else {
		orderQuery = "created_at ASC"
	}

	err := q.db.WithContext(ctx).
		Preload("Questions", func(db *gorm.DB) *gorm.DB {
			return db.Order(orderQuery)
		}).
		Preload("Questions.Answers").
		Where("id = ?", quizId).
		First(&quiz).Error

	if err != nil {
		return nil, err
	}

	return &quiz, nil
}

// FindQuizSessionByQuiz implements [IQuizSessionRepository].
func (q *QuizSessionRepositoryImpl) FindQuizSessionByQuiz(ctx context.Context) ([]models.QuizSession, error) {

	var sessions []models.QuizSession

	err := q.db.WithContext(ctx).
		Preload("Quiz").
		Preload("User").
		Order("quiz_id ASC, Quiz.created_at ASC").
		Find(&sessions).Error

	return sessions, err
}
