package classservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"giat-cerika-service/configs"
	classrequest "giat-cerika-service/internal/dto/request/class_request"
	"giat-cerika-service/internal/models"
	classrepo "giat-cerika-service/internal/repositories/class_repo"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ClassServiceImpl struct {
	classRepo classrepo.IClassRepository
	rdb       *redis.Client
}

func NewClassServiceImpl(classRepo classrepo.IClassRepository, rdb *redis.Client) IClassService {
	return &ClassServiceImpl{classRepo: classRepo, rdb: rdb}
}

func (c *ClassServiceImpl) invalidateCacheClass(ctx context.Context) {
	iter := c.rdb.Scan(ctx, 0, "classes:*", 0).Iterator()
	for iter.Next(ctx) {
		c.rdb.Del(ctx, iter.Val())
	}

	iterID := c.rdb.Scan(ctx, 0, "class:*", 0).Iterator()
	for iterID.Next(ctx) {
		c.rdb.Del(ctx, iterID.Val())
	}
}

// CreateClass implements IClassService.
func (c *ClassServiceImpl) CreateClass(ctx context.Context, req classrequest.CreateClassRequest) error {
	existClass, err := c.classRepo.FindByNameClass(ctx, req.NameClass)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get name class", 500)
	}

	if existClass != nil {
		return errorresponse.NewCustomError(errorresponse.ErrExists, "name class already exists", 409)
	}

	if strings.TrimSpace(req.NameClass) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Name Class is required", 400)
	}
	if strings.TrimSpace(req.Grade) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Grade is required", 400)
	}
	if strings.TrimSpace(req.Teacher) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Teacher is required", 400)
	}

	newClass := &models.Class{
		ID:        uuid.New(),
		NameClass: req.NameClass,
		Grade:     req.Grade,
		Teacher:   req.Teacher,
	}

	err = c.classRepo.Create(ctx, newClass)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Failed to create class", 500)
	}

	c.invalidateCacheClass(ctx)

	return nil

}

// GetAllClass implements IClassService.
func (c *ClassServiceImpl) GetAllClass(ctx context.Context, page int, limit int, search string) ([]*models.Class, int, error) {
	cacheKey := fmt.Sprintf("classes:search:%s:page:%d:limit:%d", search, page, limit)
	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var result struct {
			Data  []*models.Class `json:"data"`
			Total int             `json:"total"`
		}
		if json.Unmarshal([]byte(cached), &result) == nil {
			return result.Data, result.Total, nil
		}
	}

	offset := (page - 1) * limit

	items, total, err := c.classRepo.FindAll(ctx, limit, offset, search)
	if err != nil {
		return nil, 0, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get class", 500)
	}
	if len(items) == 0 {
		items = []*models.Class{}
	}

	buf, _ := json.Marshal(map[string]any{
		"data":  items,
		"total": total,
	})
	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)

	return items, total, nil
}

// GetByIdClass implements IClassService.
func (c *ClassServiceImpl) GetByIdClass(ctx context.Context, classId uuid.UUID) (*models.Class, error) {
	cacheKey := fmt.Sprintf("class:%s", classId)
	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var class models.Class
		if json.Unmarshal([]byte(cached), &class) == nil {
			return &class, nil
		}
	}

	class, err := c.classRepo.FindById(ctx, classId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorresponse.NewCustomError(errorresponse.ErrNotFound, "class not found", 404)
		}
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get class", 500)
	}

	buf, _ := json.Marshal(class)
	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)
	return class, nil
}

// UpdateClass implements IClassService.
func (c *ClassServiceImpl) UpdateClass(ctx context.Context, classId uuid.UUID, req classrequest.UpdateClassRequest) error {
	class, err := c.classRepo.FindById(ctx, classId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "class not found", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get class", 500)
	}

	existsClass, err := c.classRepo.FindByNameClass(ctx, req.NameClass)
	if err == nil && existsClass.ID != classId {
		return errorresponse.NewCustomError(errorresponse.ErrExists, "name class already exists", 409)
	}

	if req.NameClass != "" {
		class.NameClass = req.NameClass
	}
	if req.Grade != "" {
		class.Grade = req.Grade
	}
	if req.Teacher != "" {
		class.Teacher = req.Teacher
	}

	err = c.classRepo.Update(ctx, classId, class)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to update class", 500)
	}

	c.invalidateCacheClass(ctx)

	return nil
}

// DeleteClass implements IClassService.
func (c *ClassServiceImpl) DeleteClass(ctx context.Context, classId uuid.UUID) error {
	_, err := c.classRepo.FindById(ctx, classId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "class not found", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get class", 500)
	}

	err = c.classRepo.Delete(ctx, classId)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to delete class", 500)
	}

	c.invalidateCacheClass(ctx)

	return nil
}
