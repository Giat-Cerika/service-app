package classrepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClassRepositoryImpl struct {
	db *gorm.DB
}

func NewClassRepositoryImpl(db *gorm.DB) IClassRepository {
	return &ClassRepositoryImpl{db: db}
}

// Create implements IClassRepository.
func (c *ClassRepositoryImpl) Create(ctx context.Context, data *models.Class) error {
	return c.db.WithContext(ctx).Create(data).Error
}

// FindByNameClass implements IClassRepository.
func (c *ClassRepositoryImpl) FindByNameClass(ctx context.Context, NameClass string) (*models.Class, error) {
	var class models.Class
	if err := c.db.WithContext(ctx).First(&class, "name_class = ?", NameClass).Error; err != nil {
		return nil, err
	}

	return &class, nil
}

// FindAll implements IClassRepository.
func (c *ClassRepositoryImpl) FindAll(ctx context.Context, limit int, offset int, search string) ([]*models.Class, int, error) {
	var (
		classes []*models.Class
		count   int64
	)

	query := c.db.WithContext(ctx).Model(&models.Class{})
	if search != "" {
		query = query.Where("name_class ILIKE ?", "%"+search+"%")
	}
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&classes).Error; err != nil {
		return nil, 0, err
	}

	return classes, int(count), nil
}

// FindById implements IClassRepository.
func (c *ClassRepositoryImpl) FindById(ctx context.Context, classId uuid.UUID) (*models.Class, error) {
	var class models.Class
	if err := c.db.WithContext(ctx).First(&class, "id = ?", classId).Error; err != nil {
		return nil, err
	}

	return &class, nil
}

// Update implements IClassRepository.
func (c *ClassRepositoryImpl) Update(ctx context.Context, classId uuid.UUID, data *models.Class) error {
	return c.db.WithContext(ctx).Save(data).Error
}

// Delete implements IClassRepository.
func (c *ClassRepositoryImpl) Delete(ctx context.Context, classId uuid.UUID) error {
	return c.db.WithContext(ctx).Delete(&models.Class{}, "id = ?", classId).Error
}
