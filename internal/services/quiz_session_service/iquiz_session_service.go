package quizsessionservice

import (
	"context"
	quizrequest "giat-cerika-service/internal/dto/request/quiz_request"
	quizsessionresponse "giat-cerika-service/internal/dto/response/quiz_session_response"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type IQuizSessionService interface {
	AssignCodeQuiz(ctx context.Context, userId uuid.UUID, quizId uuid.UUID, code string) (*models.QuizSession, error)
	StartQuizSession(ctx context.Context, userId uuid.UUID, quizSessionId uuid.UUID) (*quizsessionresponse.QuizSessionStartResponse, error)
	GetQuizSessionDuration(ctx context.Context, userId uuid.UUID, quizSessionId uuid.UUID) (*quizsessionresponse.QuizSessionDurationResponse, error)
	SubmtiQuizSession(ctx context.Context, userId uuid.UUID, quizSessionId uuid.UUID, req quizrequest.SubmitQuizRequest) error
	GetOrderedQuizQuestions(ctx context.Context, userId uuid.UUID, quizSessionId uuid.UUID) (*quizsessionresponse.OrderedQuizQuestionsResponse, error)
	GetQuizSessionStudentByQuiz(ctx context.Context) ([]quizsessionresponse.ListQuestionSessionResponse, error)
}
