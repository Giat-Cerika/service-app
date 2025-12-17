package quizservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"giat-cerika-service/configs"
	quizrequest "giat-cerika-service/internal/dto/request/quiz_request"
	"giat-cerika-service/internal/models"
	quizrepo "giat-cerika-service/internal/repositories/quiz_repo"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type QuizServiceImpl struct {
	quizRepo quizrepo.IQuizRepository
	qtRepo   quizrepo.IQuizTypeRepository
	rdb      *redis.Client
}

func NewQuizServiceImpl(quizRepo quizrepo.IQuizRepository, qtRepo quizrepo.IQuizTypeRepository, rdb *redis.Client) IQuizService {
	return &QuizServiceImpl{quizRepo: quizRepo, qtRepo: qtRepo, rdb: rdb}
}

func (q *QuizServiceImpl) invalidateCacheQuiz(ctx context.Context) {
	iter := q.rdb.Scan(ctx, 0, "quizzes:*", 0).Iterator()
	for iter.Next(ctx) {
		q.rdb.Del(ctx, iter.Val())
	}

	iterID := q.rdb.Scan(ctx, 0, "quiz:*", 0).Iterator()
	for iterID.Next(ctx) {
		q.rdb.Del(ctx, iterID.Val())
	}
}

func (q *QuizServiceImpl) invalidateCacheQuestion(ctx context.Context) {
	iter := q.rdb.Scan(ctx, 0, "questions:*", 0).Iterator()
	for iter.Next(ctx) {
		q.rdb.Del(ctx, iter.Val()).Err()
	}

	IterID := q.rdb.Scan(ctx, 0, "question:*", 0).Iterator()
	for IterID.Next(ctx) {
		q.rdb.Del(ctx, IterID.Val()).Err()
	}
}

// CreateQuiz implements IQuizService.
func (q *QuizServiceImpl) CreateQuiz(ctx context.Context, req quizrequest.CreateQuizRequest) error {
	if req.QuizTypeId == uuid.Nil {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "quiz type id is required", 400)
	}
	qt, err := q.qtRepo.FindById(ctx, req.QuizTypeId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get quiz type", 500)
	}
	if strings.TrimSpace(req.Code) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "code is required", 400)
	}
	if strings.TrimSpace(req.Title) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "title is required", 400)
	}
	if strings.TrimSpace(req.Description) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "description is required", 400)
	}
	if req.StartDate.IsZero() {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "start date is required", 400)
	}
	if req.EndDate.IsZero() {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "end date is required", 400)
	}
	if req.EndDate.Before(req.StartDate) {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "end date must be after start date", 400)
	}

	newQuiz := &models.Quiz{
		ID:              uuid.New(),
		QuizTypeID:      qt.ID,
		Code:            req.Code,
		Title:           req.Title,
		Description:     req.Description,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		Status:          0,
		AmountQuestions: 0,
		AmountAssigned:  0,
	}

	err = q.quizRepo.Create(ctx, newQuiz)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to save data", 500)
	}
	q.invalidateCacheQuiz(ctx)

	return nil

}

// GetAllQuiz implements IQuizService.
func (q *QuizServiceImpl) GetAllQuiz(ctx context.Context, page int, limit int, search string) ([]*models.Quiz, int, error) {
	cacheKey := fmt.Sprintf("quizzes:search:%s:page:%d:limit:%d", search, page, limit)
	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var result struct {
			Data  []*models.Quiz `json:"data"`
			Total int            `json:"total"`
		}
		if json.Unmarshal([]byte(cached), &result) == nil {
			return result.Data, result.Total, nil
		}
	}

	offset := (page - 1) * limit

	items, total, err := q.quizRepo.FindAll(ctx, limit, offset, search)
	if err != nil {
		return nil, 0, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get quizzes", 500)
	}
	if len(items) == 0 {
		items = []*models.Quiz{}
	}
	buf, _ := json.Marshal(map[string]any{
		"data":  items,
		"total": total,
	})
	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)
	return items, total, nil
}

// GetQuizById implements IQuizService.
func (q *QuizServiceImpl) GetQuizById(ctx context.Context, quizId uuid.UUID) (*models.Quiz, error) {
	cacheKey := fmt.Sprintf("quiz:%s", quizId)
	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var quiz models.Quiz
		if json.Unmarshal([]byte(cached), &quiz) == nil {
			return &quiz, nil
		}
	}

	item, err := q.quizRepo.FindById(ctx, quizId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorresponse.NewCustomError(errorresponse.ErrNotFound, "quiz not found", 404)
		}
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get quiz", 500)
	}

	buf, _ := json.Marshal(item)
	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)
	return item, nil
}

// UpdateQuiz implements IQuizService.
func (q *QuizServiceImpl) UpdateQuiz(ctx context.Context, quizId uuid.UUID, req quizrequest.UpdateQuizRequest) error {
	quiz, err := q.quizRepo.FindById(ctx, quizId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "quiz not found", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get quiz", 500)
	}

	if req.QuizTypeId != uuid.Nil {
		quiz.QuizTypeID = req.QuizTypeId
	}
	if req.Code != "" {
		quiz.Code = req.Code
	}
	if req.Title != "" {
		quiz.Title = req.Title
	}
	if req.Description != "" {
		quiz.Description = req.Description
	}
	if !req.StartDate.IsZero() {
		quiz.StartDate = req.StartDate
	}
	if !req.EndDate.IsZero() {
		quiz.EndDate = req.EndDate
	}
	err = q.quizRepo.Update(ctx, quizId, quiz)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to update quiz", 500)
	}
	q.invalidateCacheQuiz(ctx)
	q.invalidateCacheQuestion(ctx)
	return nil
}

// DeleteQuiz implements IQuizService.
func (q *QuizServiceImpl) DeleteQuiz(ctx context.Context, quizId uuid.UUID) error {
	_, err := q.quizRepo.FindById(ctx, quizId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "quiz not found", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get quiz", 500)
	}
	err = q.quizRepo.Delete(ctx, quizId)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to delete quiz", 500)
	}
	q.invalidateCacheQuiz(ctx)
	return nil
}

// UpdateStatusQuiz implements IQuizService.
func (q *QuizServiceImpl) UpdateStatusQuiz(ctx context.Context, quizId uuid.UUID, req quizrequest.UpdateStatusQuizRequest) error {
	quiz, err := q.quizRepo.FindById(ctx, quizId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "quiz not found", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get quiz", 500)
	}
	if quiz.AmountQuestions == 0 {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "cannot activate quiz with zero questions", 400)
	}
	err = q.quizRepo.UpdateStatus(ctx, quizId, req.Status)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to update quiz status", 500)
	}
	q.invalidateCacheQuiz(ctx)
	return nil
}

// UpdateQuestionOrderMode implements [IQuizService].
func (q *QuizServiceImpl) UpdateQuestionOrderMode(ctx context.Context, quizId uuid.UUID, req quizrequest.UpdateQuestionOrderModeRequest) error {
	quiz, err := q.quizRepo.FindById(ctx, quizId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "quiz not found", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get quiz", 500)
	}
	err = q.quizRepo.UpdateQuestionOrderMode(ctx, quiz.ID, req.QuestionOrderMode)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to update question order mode", 500)
	}
	q.invalidateCacheQuiz(ctx)
	return nil
}
