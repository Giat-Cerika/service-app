package materialservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"giat-cerika-service/configs"
	rabbitmq "giat-cerika-service/configs"
	datasources "giat-cerika-service/internal/dataSources"
	materialrequest "giat-cerika-service/internal/dto/request/material_request"
	"giat-cerika-service/internal/models"
	materialrepo "giat-cerika-service/internal/repositories/material_repo"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	"giat-cerika-service/pkg/workers/payload"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
	"sync"

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

func BuildPublicID(folder, filename string) string {
	return fmt.Sprintf("%s/%s", folder, filename)
}

func Sanitize(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

func (p *MaterialServiceImpl) UploadSingle(ctx context.Context, folder, filename string, file *multipart.FileHeader) (string, string, error) {
	data, err := p.cld.UploadImage(ctx, file, folder, filename)
	if err != nil {
		return "", "", err
	}
	return data.URL, BuildPublicID(folder, filename), nil
}

func (p *MaterialServiceImpl) UploadMany(ctx context.Context, folder, prefix string, files []*multipart.FileHeader, workers int) ([]string, []string, error) {
	if len(files) == 0 {
		return nil, nil, nil
	}
	if workers < 1 {
		workers = 3
	}

	type job struct {
		idx int
		f   *multipart.FileHeader
	}
	type result struct {
		idx      int
		url      string
		publicID string
		err      error
	}

	jobs := make(chan job)
	results := make(chan result)

	worker := func() {
		for j := range jobs {
			ext := strings.ToLower(filepath.Ext(j.f.Filename))
			fileName := fmt.Sprintf("%s_%d%s", prefix, j.idx, ext)

			data, err := p.cld.UploadImage(ctx, j.f, folder, fileName)
			if err != nil {
				results <- result{idx: j.idx, err: err}
				continue
			}
			results <- result{idx: j.idx, url: data.URL, publicID: BuildPublicID(folder, fileName)}
		}
	}

	for i := 0; i < workers; i++ {
		go worker()
	}

	go func() {
		for i, f := range files {
			jobs <- job{idx: i, f: f}
		}
		close(jobs)
	}()

	urls := make([]string, len(files))
	pids := make([]string, len(files))
	var firstErr error
	done := 0
	timeout := time.After(120 * time.Second)

	for done < len(files) {
		select {
		case r := <-results:
			if r.err != nil && firstErr == nil {
				firstErr = r.err
			}
			urls[r.idx] = r.url
			pids[r.idx] = r.publicID
			done++
		case <-timeout:
		case <-ctx.Done():
		}
	}

	return urls, pids, firstErr
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

	userIDStr, ok := userID_raw.(string)
	if !ok {
		return errorresponse.NewCustomError(errorresponse.ErrUnauthorized, "invalid user id format", 401)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrUnauthorized, "invalid uuid format", 401)
	}

	if strings.TrimSpace(req.Title) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "title is required", 400)
	}
	if strings.TrimSpace(req.Description) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "description is required", 400)
	}
	if req.Cover == nil {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "cover image is required", 400)
	}
	if len(req.Gallery) == 0 {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "gallery images are required", 400)
	}

	existing, err := c.materialRepo.FindByTitle(ctx, req.Title)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to check material", 500)
	}
	if existing != nil {
		return errorresponse.NewCustomError(errorresponse.ErrExists, "material name already exists", 409)
	}

	folder := "giat-cerika-service/materials"
	nameSlug := Sanitize(req.Title)

	var wg sync.WaitGroup
	var uploadErr error

	type uploadResult struct {
		coverURL string
		gallery  []string
	}

	resultChan := make(chan uploadResult, 2)

	wg.Add(2)

	go func() {
		defer wg.Done()

		filename := "cover_" + nameSlug
		u, err := c.cld.UploadImage(ctx, req.Cover, folder, filename)
		if err != nil {
			uploadErr = err
			return
		}

		resultChan <- uploadResult{coverURL: u.URL}
	}()

	go func() {
		defer wg.Done()

		if len(req.Gallery) == 0 {
			resultChan <- uploadResult{}
			return
		}

		urls, _, err := c.UploadMany(ctx, folder, "gallery_"+nameSlug, req.Gallery, 10)
		if err != nil {
			uploadErr = err
			return
		}

		resultChan <- uploadResult{gallery: urls}
	}()

	wg.Wait()
	close(resultChan)

	if uploadErr != nil {
		return uploadErr
	}

	var coverURL string
	var galleryURLs []string

	for r := range resultChan {
		if r.coverURL != "" {
			coverURL = r.coverURL
		}
		if len(r.gallery) > 0 {
			galleryURLs = r.gallery
		}
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

		coverImg := models.Image{
			ID:        uuid.New(),
			ImagePath: coverURL,
		}

		if err := tx.Create(&coverImg).Error; err != nil {
			return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to save cover image", 500)
		}

		coverRel := models.MaterialImages{
			ID:         uuid.New(),
			MaterialID: material.ID,
			ImageID:    coverImg.ID,
			AltText:    req.Title + " Cover",
		}

		if err := tx.Create(&coverRel).Error; err != nil {
			return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to save cover relation", 500)
		}

		if len(galleryURLs) > 0 {

			images := make([]models.Image, 0, len(galleryURLs))
			for _, url := range galleryURLs {
				images = append(images, models.Image{
					ID:        uuid.New(),
					ImagePath: url,
				})
			}

			if err := tx.Create(&images).Error; err != nil {
				return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to create gallery images", 500)
			}

			relations := make([]models.MaterialImages, 0, len(images))
			for _, img := range images {
				relations = append(relations, models.MaterialImages{
					ID:         uuid.New(),
					MaterialID: material.ID,
					ImageID:    img.ID,
					AltText:    req.Title,
				})
			}

			if err := tx.Create(&relations).Error; err != nil {
				return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to create material images", 500)
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	_ = rabbitmq.PublishToQueue(
		"",
		rabbitmq.CacheInvalidateQueueName,
		payload.CacheInvalidateTask{
			Keys: []string{"materials:*"},
		},
	)

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
