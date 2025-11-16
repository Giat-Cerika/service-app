package questionnaireservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"giat-cerika-service/configs"
	questionnairerequest "giat-cerika-service/internal/dto/request/questionnaire_request"
	"giat-cerika-service/internal/models"
	questionnairerepo "giat-cerika-service/internal/repositories/questionnaire_repo"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type QuestionnaireServiceImpl struct {
	questionnaireRepo questionnairerepo.IQuestionnaireRepository
	rdb               *redis.Client
}

func NewQuestionnaireServiceImpl(questionnaireRepo questionnairerepo.IQuestionnaireRepository, rdb *redis.Client) IQuestionnaireService {
	return &QuestionnaireServiceImpl{questionnaireRepo: questionnaireRepo, rdb: rdb}
}

func (c *QuestionnaireServiceImpl) invalidateCacheQuestionnaire(ctx context.Context) {
	iter := c.rdb.Scan(ctx, 0, "questionnairees:*", 0).Iterator()
	for iter.Next(ctx) {
		c.rdb.Del(ctx, iter.Val())
	}

	iterID := c.rdb.Scan(ctx, 0, "questionnaire:*", 0).Iterator()
	for iterID.Next(ctx) {
		c.rdb.Del(ctx, iterID.Val())
	}
}

// CreateQuestionnaire implements IQuestionnaireService.
func (c *QuestionnaireServiceImpl) CreateQuestionnaire(ctx context.Context, req questionnairerequest.CreateQuestionnaireRequest) error {
	existQuestionnaire, err := c.questionnaireRepo.FindByTitle(ctx, req.Title)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get name questionnaire", 500)
	}

	if existQuestionnaire != nil {
		return errorresponse.NewCustomError(errorresponse.ErrExists, "name questionnaire already exists", 409)
	}

	if strings.TrimSpace(req.Title) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Name Questionnaire is required", 400)
	}
	if strings.TrimSpace(req.Description) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Description is required", 400)
	}
	amount, err := strconv.Atoi(req.Amount)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Amount must be a valid number", 400)
	}
	code, err := strconv.Atoi(req.Code)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Code must be a valid number", 400)
	}
	status, err := strconv.Atoi(req.Status)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Status must be a valid number", 400)
	}
	if strings.TrimSpace(req.Type) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Type is required", 400)
	}
	if strings.TrimSpace(req.Duration) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Duration is required", 400)
	}

	newQuestionnaire := &models.Questionnaire{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		Amount:      amount,
		Code:        code,
		Status:      status,
		Type:        req.Type,
		Duration:    req.Duration,
	}

	err = c.questionnaireRepo.Create(ctx, newQuestionnaire)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Failed to create questionnaire", 500)
	}

	c.invalidateCacheQuestionnaire(ctx)

	return nil

}

// GetAllQuestionnaire implements IQuestionnaireService.
func (c *QuestionnaireServiceImpl) GetAllQuestionnaire(ctx context.Context, page int, limit int, search string) ([]*models.Questionnaire, int, error) {
	cacheKey := fmt.Sprintf("questionnairees:search:%s:page:%d:limit:%d", search, page, limit)
	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var result struct {
			Data  []*models.Questionnaire `json:"data"`
			Total int                     `json:"total"`
		}
		if json.Unmarshal([]byte(cached), &result) == nil {
			return result.Data, result.Total, nil
		}
	}

	offset := (page - 1) * limit

	items, total, err := c.questionnaireRepo.FindAll(ctx, limit, offset, search)
	if err != nil {
		return nil, 0, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get questionnaire", 500)
	}
	if len(items) == 0 {
		items = []*models.Questionnaire{}
	}

	buf, _ := json.Marshal(map[string]any{
		"data":  items,
		"total": total,
	})
	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)

	return items, total, nil
}

// GetByIdQuestionnaire implements IQuestionnaireService.
func (c *QuestionnaireServiceImpl) GetByIdQuestionnaire(ctx context.Context, questionnaireId uuid.UUID) (*models.Questionnaire, error) {
	cacheKey := fmt.Sprintf("questionnaire:%s", questionnaireId)
	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var questionnaire models.Questionnaire
		if json.Unmarshal([]byte(cached), &questionnaire) == nil {
			return &questionnaire, nil
		}
	}

	questionnaire, err := c.questionnaireRepo.FindById(ctx, questionnaireId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorresponse.NewCustomError(errorresponse.ErrNotFound, "questionnaire not found", 404)
		}
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get questionnaire", 500)
	}

	buf, _ := json.Marshal(questionnaire)
	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)
	return questionnaire, nil
}

// UpdateQuestionnaire implements IQuestionnaireService.
func (c *QuestionnaireServiceImpl) UpdateQuestionnaire(ctx context.Context, questionnaireId uuid.UUID, req questionnairerequest.UpdateQuestionnaireRequest) error {
	questionnaire, err := c.questionnaireRepo.FindById(ctx, questionnaireId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "questionnaire not found", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get questionnaire", 500)
	}

	existsQuestionnaire, err := c.questionnaireRepo.FindByTitle(ctx, req.Title)
	if err == nil && existsQuestionnaire.ID != questionnaireId {
		return errorresponse.NewCustomError(errorresponse.ErrExists, "name questionnaire already exists", 409)
	}

	if req.Title != "" {
		questionnaire.Title = req.Title
	}
	if req.Description != "" {
		questionnaire.Description = req.Description
	}
	if req.Amount != "" {
		if amount, err := strconv.Atoi(req.Amount); err == nil {
			questionnaire.Amount = amount
		}
	}
	if req.Code != "" {
		if code, err := strconv.Atoi(req.Code); err == nil {
			questionnaire.Code = code
		}
	}
	if req.Status != "" {
		if status, err := strconv.Atoi(req.Status); err == nil {
			questionnaire.Status = status
		}
	}
	if req.Type != "" {
		questionnaire.Type = req.Type
	}
	if req.Duration != "" {
		questionnaire.Duration = req.Duration
	}

	err = c.questionnaireRepo.Update(ctx, questionnaireId, questionnaire)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to update questionnaire", 500)
	}

	c.invalidateCacheQuestionnaire(ctx)

	return nil
}

// DeleteQuestionnaire implements IQuestionnaireService.
func (c *QuestionnaireServiceImpl) DeleteQuestionnaire(ctx context.Context, questionnaireId uuid.UUID) error {
	_, err := c.questionnaireRepo.FindById(ctx, questionnaireId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "questionnaire not found", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get questionnaire", 500)
	}

	err = c.questionnaireRepo.Delete(ctx, questionnaireId)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to delete questionnaire", 500)
	}

	c.invalidateCacheQuestionnaire(ctx)

	return nil
}
