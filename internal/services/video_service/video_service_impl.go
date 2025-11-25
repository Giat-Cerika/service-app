package videoservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"giat-cerika-service/configs"
	videorequest "giat-cerika-service/internal/dto/request/video_request"
	"giat-cerika-service/internal/models"
	videorepo "giat-cerika-service/internal/repositories/video_repo"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type VideoServiceImpl struct {
	videoRepo videorepo.IVideoRepository
	rdb       *redis.Client
}

func NewVideoServiceImpl(videoRepo videorepo.IVideoRepository, rdb *redis.Client) IVideoService {
	return &VideoServiceImpl{videoRepo: videoRepo, rdb: rdb}
}

func (c *VideoServiceImpl) invalidateCacheVideo(ctx context.Context) {
	iter := c.rdb.Scan(ctx, 0, "videoes:*", 0).Iterator()
	for iter.Next(ctx) {
		c.rdb.Del(ctx, iter.Val())
	}

	iterID := c.rdb.Scan(ctx, 0, "video:*", 0).Iterator()
	for iterID.Next(ctx) {
		c.rdb.Del(ctx, iterID.Val())
	}
}

// CreateVideo implements IVideoService.
func (c *VideoServiceImpl) CreateVideo(
	ctx context.Context,
	req videorequest.CreateVideoRequest,
	creatorID uuid.UUID,
) error {

	existVideo, err := c.videoRepo.FindByTitle(ctx, req.Title)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get name video", 500)
	}

	if existVideo != nil {
		return errorresponse.NewCustomError(errorresponse.ErrExists, "name video already exists", 409)
	}

	if strings.TrimSpace(req.VideoPath) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Video Path is required", 400)
	}
	if strings.TrimSpace(req.Title) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Name Video is required", 400)
	}
	if strings.TrimSpace(req.Description) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Description is required", 400)
	}

	newVideo := &models.Video{
		ID:          uuid.New(),
		VideoPath:   req.VideoPath,
		Title:       req.Title,
		Description: req.Description,
		CreatedBy:   creatorID, // <-- SET CREATOR DI SINI
	}

	err = c.videoRepo.Create(ctx, newVideo)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Failed to create video", 500)
	}

	c.invalidateCacheVideo(ctx)
	return nil
}

// GetAllVideo implements IVideoService.
func (c *VideoServiceImpl) GetAllVideo(ctx context.Context, page int, limit int, search string) ([]*models.Video, int, error) {
	cacheKey := fmt.Sprintf("videoes:search:%s:page:%d:limit:%d", search, page, limit)
	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var result struct {
			Data  []*models.Video `json:"data"`
			Total int             `json:"total"`
		}
		if json.Unmarshal([]byte(cached), &result) == nil {
			return result.Data, result.Total, nil
		}
	}

	offset := (page - 1) * limit

	items, total, err := c.videoRepo.FindAll(ctx, limit, offset, search)
	if err != nil {
		return nil, 0, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get video", 500)
	}
	if len(items) == 0 {
		items = []*models.Video{}
	}

	buf, _ := json.Marshal(map[string]any{
		"data":  items,
		"total": total,
	})
	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)

	return items, total, nil
}

// GetByIdVideo implements IVideoService.
func (c *VideoServiceImpl) GetByIdVideo(ctx context.Context, videoId uuid.UUID) (*models.Video, error) {
	cacheKey := fmt.Sprintf("video:%s", videoId)
	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var video models.Video
		if json.Unmarshal([]byte(cached), &video) == nil {
			return &video, nil
		}
	}

	video, err := c.videoRepo.FindById(ctx, videoId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorresponse.NewCustomError(errorresponse.ErrNotFound, "video not found", 404)
		}
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get video", 500)
	}

	buf, _ := json.Marshal(video)
	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)
	return video, nil
}

// UpdateVideo implements IVideoService.
func (c *VideoServiceImpl) UpdateVideo(ctx context.Context, videoId uuid.UUID, req videorequest.UpdateVideoRequest) error {
	video, err := c.videoRepo.FindById(ctx, videoId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "video not found", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get video", 500)
	}

	existsVideo, err := c.videoRepo.FindByTitle(ctx, req.Title)
	if err == nil && existsVideo.ID != videoId {
		return errorresponse.NewCustomError(errorresponse.ErrExists, "name video already exists", 409)
	}

	if req.Title != "" {
		video.Title = req.Title
	}
	if req.Description != "" {
		video.Description = req.Description
	}

	err = c.videoRepo.Update(ctx, videoId, video)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to update video", 500)
	}

	c.invalidateCacheVideo(ctx)

	return nil
}

// DeleteVideo implements IVideoService.
func (c *VideoServiceImpl) DeleteVideo(ctx context.Context, videoId uuid.UUID) error {
	_, err := c.videoRepo.FindById(ctx, videoId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "video not found", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get video", 500)
	}

	err = c.videoRepo.Delete(ctx, videoId)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to delete video", 500)
	}

	c.invalidateCacheVideo(ctx)

	return nil
}

// GetAllLatestVideo implements IVideoService.
func (c *VideoServiceImpl) GetAllLatestVideo(ctx context.Context) ([]*models.Video, error) {
	cachekey := fmt.Sprintln("videoes:latest")

	if cached, err := configs.GetRedis(ctx, cachekey); err == nil && len(cached) > 0 {
		var videos []*models.Video
		if json.Unmarshal([]byte(cached), &videos) == nil {
			return videos, nil
		}
	}

	items, err := c.videoRepo.FindAllLatest(ctx)
	if err != nil {
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get latest videos", 500)
	}

	if items == nil {
		items = []*models.Video{}
	}

	buf, _ := json.Marshal(items)
	_ = configs.SetRedis(ctx, cachekey, buf, time.Minute*30)

	return items, nil
}

// GetAllPublicVideo implements IVideoService.
func (c *VideoServiceImpl) GetAllPublicVideo(ctx context.Context, page int, limit int, search string) ([]*models.Video, int, error) {
	cacheKey := fmt.Sprintf("videoes:public:search:%s:page:%d:limit:%d", search, page, limit)
	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var result struct {
			Data  []*models.Video `json:"data"`
			Total int             `json:"total"`
		}
		if json.Unmarshal([]byte(cached), &result) == nil {
			return result.Data, result.Total, nil
		}
	}

	offset := (page - 1) * limit

	items, total, err := c.videoRepo.FindAll(ctx, limit, offset, search)
	if err != nil {
		return nil, 0, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get video", 500)
	}
	if len(items) == 0 {
		items = []*models.Video{}
	}

	buf, _ := json.Marshal(map[string]any{
		"data":  items,
		"total": total,
	})
	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)

	return items, total, nil
}

// GetByIdPublicVideo implements IVideoService.
func (c *VideoServiceImpl) GetByIdPublicVideo(ctx context.Context, videoId uuid.UUID) (*models.Video, error) {
	cacheKey := fmt.Sprintf("video:public:%s", videoId)
	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var video models.Video
		if json.Unmarshal([]byte(cached), &video) == nil {
			return &video, nil
		}
	}

	video, err := c.videoRepo.FindByIdPublic(ctx, videoId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorresponse.NewCustomError(errorresponse.ErrNotFound, "video not found", 404)
		}
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get video", 500)
	}

	buf, _ := json.Marshal(video)
	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)
	return video, nil
}
