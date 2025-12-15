package quizsessionservice

import (
	"context"

	"github.com/google/uuid"
)

type IQuizSessionService interface {
	AssignCodeQuiz(ctx context.Context, userId uuid.UUID, quizId uuid.UUID, code string) (uuid.UUID, error)
}
