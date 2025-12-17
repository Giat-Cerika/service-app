package quizhistoryservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"giat-cerika-service/configs"
	"giat-cerika-service/internal/models"
	quizhistoryrepo "giat-cerika-service/internal/repositories/quiz_history_repo"
	studentrepo "giat-cerika-service/internal/repositories/student_repo"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type QuizHistoryServiceImpl struct {
	quizHistoryRepo quizhistoryrepo.IQuizHistoryRepository
	studentRepo     studentrepo.IStudentRepository
	rdb             *redis.Client
}

func NewQuizHistoryServiceImpl(quizHistoryRepo quizhistoryrepo.IQuizHistoryRepository, studentRepo studentrepo.IStudentRepository, rdb *redis.Client) IQuizHistoryService {
	return &QuizHistoryServiceImpl{quizHistoryRepo: quizHistoryRepo, studentRepo: studentRepo, rdb: rdb}
}

// GetHistoryQuizStudent implements [IQuizHistoryService].
func (q QuizHistoryServiceImpl) GetHistoryQuizStudent(ctx context.Context, userId uuid.UUID, search string) ([]*models.QuizHistory, error) {
	student, err := q.studentRepo.FindByStudentID(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorresponse.NewCustomError(errorresponse.ErrNotFound, "student not found", 404)
		}
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get student", 500)
	}

	cacheKey := fmt.Sprintf("quizHistory:%s:search:%s", userId, search)
	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var Data []*models.QuizHistory
		if json.Unmarshal([]byte(cached), &Data) == nil {
			return Data, nil
		}
	}

	items, err := q.quizHistoryRepo.FindHistoryByUserID(ctx, student.ID, search)
	if err != nil {
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get history quiz student", 500)
	}

	if len(items) == 0 {
		items = []*models.QuizHistory{}
	}

	buf, _ := json.Marshal(map[string]any{
		"data": items,
	})

	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)
	return items, nil
}
