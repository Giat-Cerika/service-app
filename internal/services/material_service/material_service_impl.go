package materialservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"giat-cerika-service/configs"
	datasources "giat-cerika-service/internal/dataSources"
	materialrequest "giat-cerika-service/internal/dto/request/material_request"
	"giat-cerika-service/internal/models"
	materialrepo "giat-cerika-service/internal/repositories/material_repo"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type MaterialServiceImpl struct {
	materialRepo materialrepo.IMaterialRepository
	rdb          *redis.Client
	cld          datasources.CloudinaryService
}

func NewMaterialServiceImpl(materialRepo materialrepo.IMaterialRepository, rdb *redis.Client, cld datasources.CloudinaryService) IMaterialService {
	return &MaterialServiceImpl{materialRepo: materialRepo, rdb: rdb, cld: cld}
}

func (c *MaterialServiceImpl) invalidateCacheMaterial(ctx context.Context) {
	iter := c.rdb.Scan(ctx, 0, "materiales:*", 0).Iterator()
	for iter.Next(ctx) {
		c.rdb.Del(ctx, iter.Val())
	}

	iterID := c.rdb.Scan(ctx, 0, "material:*", 0).Iterator()
	for iterID.Next(ctx) {
		c.rdb.Del(ctx, iterID.Val())
	}
}

// CreateMaterial implements IMaterialService.
func (c *MaterialServiceImpl) CreateMaterial(ctx context.Context, req materialrequest.CreateMaterialRequest) error {
    userID_raw := ctx.Value("user_id")
    if userID_raw == nil {
        return errorresponse.NewCustomError(errorresponse.ErrUnauthorized, "user not authorized", 401)
    }
    userID, ok := userID_raw.(uuid.UUID)
    if !ok {
        return errorresponse.NewCustomError(errorresponse.ErrUnauthorized, "invalid user id", 401)
    }

    if strings.TrimSpace(req.Title) == "" {
        return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "name material is required", 400)
    }
    if strings.TrimSpace(req.Description) == "" {
        return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "description is required", 400)
    }
    if req.Cover == nil {
        return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Cover image is required", 400)
    }
    if len(req.Gallery) == 0 {
        return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "images are required", 400)
    }

    existing, err := c.materialRepo.FindByTitle(ctx, req.Title)
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get material", 500)
    }
    if existing != nil {
        return errorresponse.NewCustomError(errorresponse.ErrExists, "material name already exists", 409)
    }

    err = configs.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

        material := &models.Materials{
            ID:          uuid.New(),
            Title:       req.Title,
            Description: req.Description,
            CreatedBy:   userID,  
        }

        if err := tx.Create(material).Error; err != nil {
            return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to create material", 500)
        }

        return nil
    })
    if err != nil {
        return err
    }

    c.invalidateCacheMaterial(ctx)
    return nil
}

// GetAllMaterial implements IMaterialService.
func (c *MaterialServiceImpl) GetAllMaterial(ctx context.Context, page int, limit int, search string) ([]*models.Materials, int, error) {
	cacheKey := fmt.Sprintf("materiales:search:%s:page:%d:limit:%d", search, page, limit)
	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var result struct {
			Data  []*models.Materials `json:"data"`
			Total int                 `json:"total"`
		}
		if json.Unmarshal([]byte(cached), &result) == nil {
			return result.Data, result.Total, nil
		}
	}

	offset := (page - 1) * limit

	items, total, err := c.materialRepo.FindAll(ctx, limit, offset, search)
	if err != nil {
		return nil, 0, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get material", 500)
	}
	if len(items) == 0 {
		items = []*models.Materials{}
	}

	buf, _ := json.Marshal(map[string]any{
		"data":  items,
		"total": total,
	})
	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)

	return items, total, nil
}

// GetByIdMaterial implements IMaterialService.
func (c *MaterialServiceImpl) GetByIdMaterial(ctx context.Context, materialId uuid.UUID) (*models.Materials, error) {
	cacheKey := fmt.Sprintf("material:%s", materialId)
	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var material models.Materials
		if json.Unmarshal([]byte(cached), &material) == nil {
			return &material, nil
		}
	}

	material, err := c.materialRepo.FindById(ctx, materialId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorresponse.NewCustomError(errorresponse.ErrNotFound, "material not found", 404)
		}
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get material", 500)
	}

	buf, _ := json.Marshal(material)
	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)
	return material, nil
}

// UpdateMaterial implements IMaterialService.
func (c *MaterialServiceImpl) UpdateMaterial(ctx context.Context, materialId uuid.UUID, req materialrequest.UpdateMaterialRequest) error {
	material, err := c.materialRepo.FindById(ctx, materialId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "material not found", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get material", 500)
	}

	existsMaterial, err := c.materialRepo.FindByTitle(ctx, req.Title)
	if err == nil && existsMaterial.ID != materialId {
		return errorresponse.NewCustomError(errorresponse.ErrExists, "name material already exists", 409)
	}

	if req.Title != "" {
		material.Title = req.Title
	}
	if req.Description != "" {
		material.Description = req.Description
	}

	err = c.materialRepo.Update(ctx, materialId, material)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to update material", 500)
	}

	c.invalidateCacheMaterial(ctx)

	return nil
}

// DeleteMaterial implements IMaterialService.
func (c *MaterialServiceImpl) DeleteMaterial(ctx context.Context, materialId uuid.UUID) error {
	_, err := c.materialRepo.FindById(ctx, materialId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "material not found", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get material", 500)
	}

	err = c.materialRepo.Delete(ctx, materialId)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to delete material", 500)
	}

	c.invalidateCacheMaterial(ctx)

	return nil
}
