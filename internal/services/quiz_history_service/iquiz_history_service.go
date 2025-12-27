package quizhistoryservice

import (
	"context"
	quizhistoryresponse "giat-cerika-service/internal/dto/response/quiz_history_response"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type IQuizHistoryService interface {
	GetHistoryQuizStudent(ctx context.Context, userId uuid.UUID, search string) ([]quizhistoryresponse.QuizHistoryResponse, error)
	GetAllHistoryQuestionByQuizHistory(ctx context.Context, quizHistoryId uuid.UUID) ([]*models.QuestionHistory, error)
	GetHistoryQuizByQuizID(ctx context.Context) ([]quizhistoryresponse.QuizHistoryGroupAdminResponse, error)
}
