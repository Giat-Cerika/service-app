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

type QuizTypeServiceImpl struct {
	qtRepo quizrepo.IQuizTypeRepository
	rdb    *redis.Client
}

func NewQuizTypeServiceImpl(qtRepo quizrepo.IQuizTypeRepository, rdb *redis.Client) IQuizTypeService {
	return &QuizTypeServiceImpl{qtRepo: qtRepo, rdb: rdb}
}

func (q *QuizTypeServiceImpl) invalideCacheQT(ctx context.Context) {
	iter := q.rdb.Scan(ctx, 0, "quizTypes:*", 0).Iterator()
	for iter.Next(ctx) {
		q.rdb.Del(ctx, iter.Val())
	}

	iterID := q.rdb.Scan(ctx, 0, "quizType:*", 0).Iterator()
	for iterID.Next(ctx) {
		q.rdb.Del(ctx, iterID.Val())
	}
}

// CreateQt implements IQuizTypeService.
func (q *QuizTypeServiceImpl) CreateQt(ctx context.Context, req quizrequest.CreateQuizTypeRequest) error {
	existQt, err := q.qtRepo.FindByName(ctx, req.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get name quiz type", 500)
	}

	if existQt != nil {
		return errorresponse.NewCustomError(errorresponse.ErrExists, "name quiz type already exists", 409)
	}

	if strings.TrimSpace(req.Name) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "name quiz type is required", 400)
	}

	if strings.TrimSpace(req.Description) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "description quiz type is required", 400)
	}

	newQt := &models.QuizType{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
	}

	err = q.qtRepo.Create(ctx, newQt)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to save data", 500)
	}

	q.invalideCacheQT(ctx)

	return nil

}

// GetAllQt implements IQuizTypeService.
func (q *QuizTypeServiceImpl) GetAllQt(ctx context.Context) ([]*models.QuizType, error) {
	cacheKey := fmt.Sprintln("quizTypes:all")
	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var qt []*models.QuizType
		if json.Unmarshal([]byte(cached), &qt) == nil {
			return qt, nil
		}
	}

	items, err := q.qtRepo.FindAll(ctx)
	if err != nil {
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get quiz type", 500)
	}

	if len(items) == 0 {
		items = []*models.QuizType{}
	}

	buf, _ := json.Marshal(items)
	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)
	return items, nil
}

// GetByIdQt implements IQuizTypeService.
func (q *QuizTypeServiceImpl) GetByIdQt(ctx context.Context, quizTypeId uuid.UUID) (*models.QuizType, error) {
	cacheKey := fmt.Sprintf("quizType:%s", quizTypeId.String())
	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var qt models.QuizType
		if json.Unmarshal([]byte(cached), &qt) == nil {
			return &qt, nil
		}
	}
	item, err := q.qtRepo.FindById(ctx, quizTypeId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorresponse.NewCustomError(errorresponse.ErrNotFound, "quiz type not found", 404)
		}
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get quiz type", 500)
	}
	buf, _ := json.Marshal(item)
	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)
	return item, nil
}

// UpdateQt implements IQuizTypeService.
func (q *QuizTypeServiceImpl) UpdateQt(ctx context.Context, quizTypeId uuid.UUID, req quizrequest.UpdateQuizTypeRequest) error {
	qt, err := q.qtRepo.FindById(ctx, quizTypeId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "quiz type not found", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get quiz type", 500)
	}
	if req.Name != "" {
		qt.Name = req.Name
	}
	if req.Description != "" {
		qt.Description = req.Description
	}
	err = q.qtRepo.Update(ctx, quizTypeId, qt)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to update quiz type", 500)
	}
	q.invalideCacheQT(ctx)

	return nil
}

// DeleteQt implements IQuizTypeService.
func (q *QuizTypeServiceImpl) DeleteQt(ctx context.Context, quizTypeId uuid.UUID) error {
	_, err := q.qtRepo.FindById(ctx, quizTypeId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "quiz type not found", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get quiz type", 500)
	}
	err = q.qtRepo.Delete(ctx, quizTypeId)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to delete quiz type", 500)
	}
	q.invalideCacheQT(ctx)

	return nil
}
