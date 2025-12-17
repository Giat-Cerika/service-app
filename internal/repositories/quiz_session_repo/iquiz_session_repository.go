package quizsessionrepo

import (
	"context"
	"giat-cerika-service/internal/models"
	"time"

	"github.com/google/uuid"
)

type IQuizSessionRepository interface {
	AssignCodeQuiz(ctx context.Context, quizId uuid.UUID, code string) (uuid.UUID, error)
	SaveQuizSession(ctx context.Context, data *models.QuizSession) error
	CreateStartedAt(ctx context.Context, quizSessionId uuid.UUID) error
	FindById(ctx context.Context, userId uuid.UUID, quizSessionId uuid.UUID) (*models.QuizSession, error)
	FindByUserAndQuiz(ctx context.Context, userId uuid.UUID, quizId uuid.UUID) (*models.QuizSession, error)

	BulkSaveResponses(ctx context.Context, quizSessionId uuid.UUID, data []*models.Response) error
	CompleteQuizSession(ctx context.Context, quizSessionId uuid.UUID, score int, maxScore int, completedAt *time.Time) error
	SaveQuizHistoryTransaction(
		ctx context.Context,
		quizHistory *models.QuizHistory,
		questionHistories []models.QuestionHistory,
		answerHistories []models.AnswerHistory,
	) error
	FindQuizWithOrderedQuestions(ctx context.Context, quizId uuid.UUID, orderMode string) (*models.Quiz, error)

	FindQuizSessionByQuiz(ctx context.Context) ([]models.QuizSession, error)
}
