package quizsessionservice

import (
	"context"
	"errors"
	"giat-cerika-service/internal/models"
	quizrepo "giat-cerika-service/internal/repositories/quiz_repo"
	quizsessionrepo "giat-cerika-service/internal/repositories/quiz_session_repo"
	studentrepo "giat-cerika-service/internal/repositories/student_repo"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuizSessionServiceImpl struct {
	quizSessionRepo quizsessionrepo.IQuizSessionRepository
	quizRepo        quizrepo.IQuizRepository
	studentRepo     studentrepo.IStudentRepository
}

func NewQuizSessionServiceImpl(qsRepo quizsessionrepo.IQuizSessionRepository, quizRepo quizrepo.IQuizRepository, studentRepo studentrepo.IStudentRepository) IQuizSessionService {
	return &QuizSessionServiceImpl{quizSessionRepo: qsRepo, quizRepo: quizRepo, studentRepo: studentRepo}
}

// AssignCodeQuiz implements [IQuizSessionService].
func (q *QuizSessionServiceImpl) AssignCodeQuiz(ctx context.Context, userId uuid.UUID, quizId uuid.UUID, code string) (uuid.UUID, error) {
	quiz, err := q.quizRepo.FindById(ctx, quizId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return uuid.Nil, errorresponse.NewCustomError(errorresponse.ErrNotFound, "quiz not found", 404)
		}
		return uuid.Nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get quiz", 500)
	}

	student, err := q.studentRepo.FindByStudentID(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return uuid.Nil, errorresponse.NewCustomError(errorresponse.ErrNotFound, "student not found", 404)
		}
		return uuid.Nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get student", 500)
	}

	if strings.TrimSpace(code) == "" {
		return uuid.Nil, errorresponse.NewCustomError(errorresponse.ErrBadRequest, "code is required", 400)
	}

	quizId, err = q.quizSessionRepo.AssignCodeQuiz(ctx, quiz.ID, code)
	if err != nil {
		return uuid.Nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get quiz data", 500)
	}

	if quizId != quiz.ID || code != quiz.Code {
		return uuid.Nil, errorresponse.NewCustomError(errorresponse.ErrExists, "code access doesn't matching", 409)
	}

	if quizId == quiz.ID && quiz.Status == 0 {
		return uuid.Nil, errorresponse.NewCustomError(errorresponse.ErrBadRequest, "quiz has been drafts", 400)
	}

	newQuizSession := &models.QuizSession{
		ID:       uuid.New(),
		UserID:   student.ID,
		QuizID:   quiz.ID,
		Score:    0,
		MaxScore: 0,
		Status:   models.SessionStatusStarted,
	}

	if err := q.quizSessionRepo.SaveQuizSession(ctx, newQuizSession); err != nil {
		return uuid.Nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to save quiz session", 500)
	}

	return quiz.ID, nil
}
