package rolerepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleRepositoryImpl struct {
	db *gorm.DB
}

func NewRoleRepositoryImpl(db *gorm.DB) IRoleRepository {
	return &RoleRepositoryImpl{db: db}
}

// Create implements IRoleRepository.
func (r *RoleRepositoryImpl) Create(ctx context.Context, data *models.Role) error {
	return r.db.WithContext(ctx).Create(data).Error
}

// FindByName implements IRoleRepository.
func (r *RoleRepositoryImpl) FindByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role

	if err := r.db.WithContext(ctx).First(&role, "name = ?", name).Error; err != nil {
		return nil, err
	}

	return &role, nil
}

// FindAll implements IRoleRepository.
func (r *RoleRepositoryImpl) FindAll(ctx context.Context, limit int, offset int, search string) ([]*models.Role, int, error) {
	var (
		role  []*models.Role
		count int64
	)

	query := r.db.WithContext(ctx).Model(&models.Role{})
	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&role).Error; err != nil {
		return nil, 0, err
	}

	return role, int(count), nil
}

// FindById implements IRoleRepository.
func (r *RoleRepositoryImpl) FindById(ctx context.Context, roleId uuid.UUID) (*models.Role, error) {
	var role models.Role
	if err := r.db.WithContext(ctx).First(&role, "id = ?", roleId).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// Update implements IRoleRepository.
func (r *RoleRepositoryImpl) Update(ctx context.Context, roleId uuid.UUID, data *models.Role) error {
	return r.db.WithContext(ctx).Save(data).Error
}

// Delete implements IRoleRepository.
func (r *RoleRepositoryImpl) Delete(ctx context.Context, roleId uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Role{}, "id = ?", roleId).Error
}
