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
	adminrepo "giat-cerika-service/internal/repositories/admin_repo"
	materialrepo "giat-cerika-service/internal/repositories/material_repo"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	rabbitmq "giat-cerika-service/pkg/constant/rabbitMq"
	"giat-cerika-service/pkg/utils"
	"giat-cerika-service/pkg/workers/payload"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type MaterialServiceImpl struct {
	materialRepo materialrepo.IMaterialRepository
	adminRepo    adminrepo.IAdminRepository
	rdb          *redis.Client
	cld          datasources.CloudinaryService
}

func NewMaterialServiceImpl(materialRepo materialrepo.IMaterialRepository, adminRepo adminrepo.IAdminRepository, rdb *redis.Client, cld datasources.CloudinaryService) IMaterialService {
	return &MaterialServiceImpl{materialRepo: materialRepo, adminRepo: adminRepo, rdb: rdb, cld: cld}
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

func fileMateriToBytes(fh *multipart.FileHeader) ([]byte, error) {
	f, err := fh.Open()
	if err != nil {
		return nil, err
	}

	defer f.Close()
	return io.ReadAll(f)
}

var PublishImageAsync = func(p payload.ImageUploadPayload) {
	go func() {
		_ = rabbitmq.PublishToQueue(
			"",
			rabbitmq.SendImageMateriQueueName,
			p,
		)
	}()
}

// CreateMaterial implements IMaterialService.
func (c *MaterialServiceImpl) CreateMaterial(ctx context.Context, adminId uuid.UUID, req materialrequest.CreateMaterialRequest) error {
	admin, err := c.adminRepo.FindAdmin(ctx, adminId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "users need access permission", 401)
	}

	if strings.TrimSpace(req.Title) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "title material is required", 400)
	}
	if strings.TrimSpace(req.Description) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "description is required", 400)
	}
	if req.Cover == nil {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Cover image is required", 400)
	}
	if req.Gallery == nil {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "galerry image is required", 400)
	}
	if len(req.Gallery) == 0 {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "gallery are required", 400)
	}

	existing, err := c.materialRepo.FindByTitle(ctx, req.Title)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get material", 500)
	}
	if existing != nil {
		return errorresponse.NewCustomError(errorresponse.ErrExists, "material name already exists", 409)
	}

	materi := &models.Materials{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		CreatedBy:   admin.ID,
	}

	if err := c.materialRepo.Create(ctx, materi); err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to create materi", 500)
	}

	if req.Cover != nil {
		if bin, err := fileMateriToBytes(req.Cover); err == nil && len(bin) > 0 {
			go PublishImageAsync(payload.ImageUploadPayload{
				ID:        materi.ID,
				Type:      "single",
				FileBytes: bin,
				Folder:    "giat_ceria/materials",
				Filename:  fmt.Sprintf("materi_%s_cover", materi.ID.String()),
			})
		}
	}
	for i, g := range req.Gallery {
		if g == nil {
			continue
		}
		if bin, err := fileMateriToBytes(g); err == nil && len(bin) > 0 {
			go PublishImageAsync(payload.ImageUploadPayload{
				ID:        materi.ID,
				Type:      "many",
				FileBytes: bin,
				Folder:    "giat_ceria/materials",
				Filename:  fmt.Sprintf("materi_%s_gallery_%d", materi.ID, i+1),
			})

		}
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

	// cover handling (sama seperti sebelumnya)
	if req.Cover != nil {
		if material.Cover != "" {
			publicID := utils.ExtractPublicIDFromCloudinaryURL(material.Cover)
			if publicID != "" {
				_ = c.cld.DestroyImage(ctx, publicID)
			}
		}
		if bin, err := fileMateriToBytes(req.Cover); err == nil && len(bin) > 0 {
			// PublishImageAsync sudah membuat goroutine sendiri => panggil langsung
			PublishImageAsync(payload.ImageUploadPayload{
				ID:        material.ID,
				Type:      "single",
				FileBytes: bin,
				Folder:    "giat_ceria/materials",
				Filename:  fmt.Sprintf("materi_%s_cover", material.ID.String()),
			})
		}
	}

	if len(req.Gallery) > 0 {
		if req.ReplaceGallery {
			// Hapus gallery di DB dulu, tangani error
			if err := c.materialRepo.DeleteGalleryByMateriId(ctx, material.ID); err != nil {
				// kalau gagal delete, beri error agar tidak inconsistent
				return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to delete existing gallery", 500)
			}

			// Hapus file di cloudinary berdasarkan data yang ada di memory (material.MaterialImages)
			// gunakan copy untuk safety
			for _, gi := range material.MaterialImages {
				if gi.Image.ImagePath != "" {
					publicID := utils.ExtractPublicIDFromCloudinaryURL(gi.Image.ImagePath)
					if publicID != "" {
						_ = c.cld.DestroyImage(ctx, publicID) // gagal destroy tidak fatal
					}
				}
			}

			// Sangat penting: kosongkan slice agar state in-memory sesuai dengan DB
			material.MaterialImages = []models.MaterialImages{}
		}

		// upload file baru (buat nama file lebih unik untuk menghindari collision)
		for i, g := range req.Gallery {
			if g == nil {
				continue
			}
			if bin, err := fileMateriToBytes(g); err == nil && len(bin) > 0 {
				go PublishImageAsync(payload.ImageUploadPayload{
					ID:        material.ID,
					Type:      "many",
					FileBytes: bin,
					Folder:    "giat_ceria/materials",
					// gunakan String() + index + timestamp supaya unik
					Filename: fmt.Sprintf("materi_%s_gallery_%d_%d", material.ID.String(), i+1, time.Now().UnixNano()),
				})
			}
		}
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
	materi, err := c.materialRepo.FindById(ctx, materialId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "material not found", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get material", 500)
	}

	if len(materi.MaterialImages) > 0 {
		_ = c.materialRepo.DeleteGalleryByMateriId(ctx, materi.ID)
	}

	if materi.Cover != "" {
		publicID := utils.ExtractPublicIDFromCloudinaryURL(materi.Cover)
		if publicID != "" {
			if err := c.cld.DestroyImage(ctx, publicID); err != nil {
				return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to delete image", 500)
			}
		}
	}

	err = c.materialRepo.Delete(ctx, materialId)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to delete material", 500)
	}

	c.invalidateCacheMaterial(ctx)

	return nil
}

// GetAllLatestMaterial implements IMaterialService.
func (c *MaterialServiceImpl) GetAllLatestMaterial(ctx context.Context) ([]*models.Materials, error) {
	cacheKey := fmt.Sprintln("materiales:latest")

	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var materiales []*models.Materials
		if json.Unmarshal([]byte(cached), &materiales) == nil {
			return materiales, nil
		}
	}

	items, err := c.materialRepo.FindAllLatest(ctx)
	if err != nil {
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get latest material", 500)
	}

	if items == nil {
		items = []*models.Materials{}
	}

	buf, _ := json.Marshal(items)
	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)

	return items, nil
}

// GetAllPublicMaterial implements IMaterialService.
func (c *MaterialServiceImpl) GetAllPublicMaterial(ctx context.Context, page int, limit int, search string) ([]*models.Materials, int, error) {
	cacheKey := fmt.Sprintf("materiales:public:search:%s:page:%d:limit:%d", search, page, limit)
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

// GetByIdPublicMaterial implements IMaterialService.
func (c *MaterialServiceImpl) GetByIdPublicMaterial(ctx context.Context, materialId uuid.UUID) (*models.Materials, error) {
	cacheKey := fmt.Sprintf("material:public:%s", materialId)
	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var material models.Materials
		if json.Unmarshal([]byte(cached), &material) == nil {
			return &material, nil
		}
	}

	material, err := c.materialRepo.FindByIdPublic(ctx, materialId)
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
