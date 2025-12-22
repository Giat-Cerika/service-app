package quizhistoryservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"giat-cerika-service/configs"
	quizhistoryresponse "giat-cerika-service/internal/dto/response/quiz_history_response"
	"giat-cerika-service/internal/models"
	quizhistoryrepo "giat-cerika-service/internal/repositories/quiz_history_repo"
	quizrepo "giat-cerika-service/internal/repositories/quiz_repo"
	studentrepo "giat-cerika-service/internal/repositories/student_repo"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	"giat-cerika-service/pkg/utils"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type QuizHistoryServiceImpl struct {
	quizHistoryRepo quizhistoryrepo.IQuizHistoryRepository
	studentRepo     studentrepo.IStudentRepository
	quizRepo        quizrepo.IQuizRepository
	rdb             *redis.Client
}

func NewQuizHistoryServiceImpl(quizHistoryRepo quizhistoryrepo.IQuizHistoryRepository, studentRepo studentrepo.IStudentRepository, quizRepo quizrepo.IQuizRepository, rdb *redis.Client) IQuizHistoryService {
	return &QuizHistoryServiceImpl{quizHistoryRepo: quizHistoryRepo, studentRepo: studentRepo, quizRepo: quizRepo, rdb: rdb}
}

// GetHistoryQuizStudent implements [IQuizHistoryService].
func (q QuizHistoryServiceImpl) GetHistoryQuizStudent(
	ctx context.Context,
	userId uuid.UUID,
	search string,
) ([]quizhistoryresponse.QuizHistoryResponse, error) {

	student, err := q.studentRepo.FindByStudentID(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorresponse.NewCustomError(errorresponse.ErrNotFound, "student not found", 404)
		}
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get student", 500)
	}

	cacheKey := fmt.Sprintf("quizHistory:%s:search:%s", userId, search)
	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var data []quizhistoryresponse.QuizHistoryResponse
		if json.Unmarshal([]byte(cached), &data) == nil {
			return data, nil
		}
	}

	items, err := q.quizHistoryRepo.FindHistoryByUserID(ctx, student.ID, search)
	if err != nil {
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get history quiz student", 500)
	}

	res := make([]quizhistoryresponse.QuizHistoryResponse, 0, len(items))

	for _, h := range items {
		// snapshot
		currentAssigned := h.AmountAssigned

		// kalau quiz masih ada → pakai global
		if quiz, err := q.quizRepo.FindById(ctx, h.QuizID); err == nil {
			currentAssigned = quiz.AmountAssigned
		}

		res = append(res, quizhistoryresponse.QuizHistoryResponse{
			ID:              h.ID,
			Title:           h.Title,
			Description:     h.Description,
			StartDate:       utils.FormatDateTime(h.StartDate),
			EndDate:         utils.FormatDateTime(h.EndDate),
			AmountQuestions: h.AmountQuestions,
			AmountAssigned:  currentAssigned,
			Score:           h.Score,
			MaxScore:        h.MaxScore,
			Percentage:      h.Percentage,
			Status:          string(h.Status),
			StartedAt:       utils.FormatDateTime(h.StartedAt),
			CompletedAt:     utils.FormatDateTime(h.CompletedAt),
			StatusCategory:  h.StatusCategory,
			CreatedAt:       utils.FormatDate(h.CreatedAt),
			UpdatedAt:       utils.FormatDate(h.UpdatedAt),
		})
	}

	buf, _ := json.Marshal(map[string]any{
		"data": res,
	})
	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)

	return res, nil
}

// GetAllHistoryQuestionByQuizHistory implements [IQuizHistoryService].
func (q *QuizHistoryServiceImpl) GetAllHistoryQuestionByQuizHistory(ctx context.Context, quizHistoryId uuid.UUID) ([]*models.QuestionHistory, error) {
	cacheKey := fmt.Sprintf("questions_history:quiz_history:%s", quizHistoryId)

	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var data []*models.QuestionHistory
		if json.Unmarshal([]byte(cached), &data) == nil {
			return data, nil
		}
	}

	quizHistory, err := q.quizHistoryRepo.FindQuizHistoryById(ctx, quizHistoryId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorresponse.NewCustomError(errorresponse.ErrNotFound, "quiz history not found", 404)
		}
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get quiz history", 500)
	}

	items, err := q.quizHistoryRepo.FindAllQuestionHistory(ctx, quizHistory.ID)
	if err != nil {
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get question history", 500)
	}

	if items == nil {
		items = []*models.QuestionHistory{}
	}

	buf, _ := json.Marshal(map[string]any{
		"data": items,
	})

	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)

	return items, nil
}
