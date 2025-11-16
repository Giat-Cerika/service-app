package videorepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VideoRepositoryImpl struct {
	db *gorm.DB
}

func NewVideoRepositoryImpl(db *gorm.DB) IVideoRepository {
	return &VideoRepositoryImpl{db: db}
}

// Create implements IVideoRepository.
func (c *VideoRepositoryImpl) Create(ctx context.Context, data *models.Video) error {
	return c.db.WithContext(ctx).Create(data).Error
}

// FindByTitle implements IVideoRepository.
func (c *VideoRepositoryImpl) FindByTitle(ctx context.Context, Title string) (*models.Video, error) {
	var video models.Video
	if err := c.db.WithContext(ctx).First(&video, "title = ?", Title).Error; err != nil {
		return nil, err
	}

	return &video, nil
}

// FindAll implements IVideoRepository.
func (c *VideoRepositoryImpl) FindAll(ctx context.Context, limit int, offset int, search string) ([]*models.Video, int, error) {
	var (
		videoes []*models.Video
		count   int64
	)

	query := c.db.WithContext(ctx).Model(&models.Video{})
	if search != "" {
		query = query.Where("name_video ILIKE ?", "%"+search+"%")
	}
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&videoes).Error; err != nil {
		return nil, 0, err
	}

	return videoes, int(count), nil
}

// FindById implements IVideoRepository.
func (c *VideoRepositoryImpl) FindById(ctx context.Context, videoId uuid.UUID) (*models.Video, error) {
	var video models.Video
	if err := c.db.WithContext(ctx).First(&video, "id = ?", videoId).Error; err != nil {
		return nil, err
	}

	return &video, nil
}

// Update implements IVideoRepository.
func (c *VideoRepositoryImpl) Update(ctx context.Context, videoId uuid.UUID, data *models.Video) error {
	return c.db.WithContext(ctx).Save(data).Error
}

// Delete implements IVideoRepository.
func (c *VideoRepositoryImpl) Delete(ctx context.Context, videoId uuid.UUID) error {
	return c.db.WithContext(ctx).Delete(&models.Video{}, "id = ?", videoId).Error
}
