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

func (c *MaterialRepositoryImpl) preloadRelations(db *gorm.DB) *gorm.DB {
	return db.
		Preload("MaterialImages.Image")
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
		count      int64
	)

	query := c.db.WithContext(ctx).Model(&models.Materials{})
	if search != "" {
		query = query.Where("title ILIKE ?", "%"+search+"%")
	}
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	query = c.preloadRelations(query)
	if err := query.Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&materiales).Error; err != nil {
		return nil, 0, err
	}

	return materiales, int(count), nil
}

// FindById implements IMaterialRepository.
func (c *MaterialRepositoryImpl) FindById(ctx context.Context, materialId uuid.UUID) (*models.Materials, error) {
	var material models.Materials
	if err := c.preloadRelations(c.db.WithContext(ctx)).First(&material, "id = ?", materialId).Error; err != nil {
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

// UpdateCoverMateri implements IMaterialRepository.
func (c *MaterialRepositoryImpl) UpdateCoverMateri(ctx context.Context, materiId uuid.UUID, cover string) error {
	return c.db.WithContext(ctx).Model(&models.Materials{}).
		Where("id = ?", materiId).
		Update("cover", cover).Error
}

// CreateGallery implements IMaterialRepository.
func (c *MaterialRepositoryImpl) CreateGallery(ctx context.Context, materiId uuid.UUID, imageId uuid.UUID, alt string) error {
	materiImage := &models.MaterialImages{
		ID:         uuid.New(),
		MaterialID: materiId,
		ImageID:    imageId,
		AltText:    alt,
	}

	return c.db.WithContext(ctx).Create(materiImage).Error
}

// CreateImage implements IMaterialRepository.
func (c *MaterialRepositoryImpl) CreateImage(ctx context.Context, image *models.Image) error {
	return c.db.WithContext(ctx).Create(image).Error
}

// DeleteGalleryByMateriId implements IMaterialRepository.
func (c *MaterialRepositoryImpl) DeleteGalleryByMateriId(ctx context.Context, materiId uuid.UUID) error {
	return c.db.WithContext(ctx).
		Where("material_id = ?", materiId).
		Delete(&models.MaterialImages{}).Error
}

// FindAllLatest implements IMaterialRepository.
func (c *MaterialRepositoryImpl) FindAllLatest(ctx context.Context) ([]*models.Materials, error) {
	var materials []*models.Materials
	query := c.db.WithContext(ctx).Model(&models.Materials{})
	query = c.preloadRelations(query)
	if err := query.
		Order("created_at DESC"). // Urutkan dari yang terbaru
		Limit(5).                 // Ambil 5 data teratas
		Find(&materials).Error; err != nil {
		return nil, err
	}
	return materials, nil
}

// FindAllPublic implements IMaterialRepository.
func (c *MaterialRepositoryImpl) FindAllPublic(ctx context.Context, limit int, offset int, search string) ([]*models.Materials, int, error) {
	var (
		materiales []*models.Materials
		count      int64
	)

	query := c.db.WithContext(ctx).Model(&models.Materials{})
	if search != "" {
		query = query.Where("title ILIKE ?", "%"+search+"%")
	}
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	query = c.preloadRelations(query)
	if err := query.Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&materiales).Error; err != nil {
		return nil, 0, err
	}

	return materiales, int(count), nil
}

// FindByIdPublic implements IMaterialRepository.
func (c *MaterialRepositoryImpl) FindByIdPublic(ctx context.Context, materialId uuid.UUID) (*models.Materials, error) {
	var material models.Materials
	if err := c.preloadRelations(c.db.WithContext(ctx)).First(&material, "id = ?", materialId).Error; err != nil {
		return nil, err
	}

	return &material, nil
}
