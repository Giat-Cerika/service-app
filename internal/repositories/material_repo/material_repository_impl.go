package materialrepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MaterialRepositoryImpl struct {
	db *gorm.DB
}

func NewMaterialRepositoryImpl(db *gorm.DB) IMaterialRepository {
	return &MaterialRepositoryImpl{db: db}
}

// Create implements IMaterialRepository.
func (c *MaterialRepositoryImpl) Create(ctx context.Context, data *models.Materials) error {
	return c.db.WithContext(ctx).Create(data).Error
}

// FindByTitle implements IMaterialRepository.
func (c *MaterialRepositoryImpl) FindByTitle(ctx context.Context, Title string) (*models.Materials, error) {
	var material models.Materials
	if err := c.db.WithContext(ctx).First(&material, "title = ?", Title).Error; err != nil {
		return nil, err
	}

	return &material, nil
}

// FindAll implements IMaterialRepository.
func (c *MaterialRepositoryImpl) FindAll(ctx context.Context, limit int, offset int, search string) ([]*models.Materials, int, error) {
	var (
		materiales []*models.Materials
		count   int64
	)

	query := c.db.WithContext(ctx).Model(&models.Materials{})
	if search != "" {
		query = query.Where("name_material ILIKE ?", "%"+search+"%")
	}
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&materiales).Error; err != nil {
		return nil, 0, err
	}

	return materiales, int(count), nil
}

// FindById implements IMaterialRepository.
func (c *MaterialRepositoryImpl) FindById(ctx context.Context, materialId uuid.UUID) (*models.Materials, error) {
	var material models.Materials
	if err := c.db.WithContext(ctx).First(&material, "id = ?", materialId).Error; err != nil {
		return nil, err
	}

	return &material, nil
}

// Update implements IMaterialRepository.
func (c *MaterialRepositoryImpl) Update(ctx context.Context, materialId uuid.UUID, data *models.Materials) error {
	return c.db.WithContext(ctx).Save(data).Error
}

// Delete implements IMaterialRepository.
func (c *MaterialRepositoryImpl) Delete(ctx context.Context, materialId uuid.UUID) error {
	return c.db.WithContext(ctx).Delete(&models.Materials{}, "id = ?", materialId).Error
}
