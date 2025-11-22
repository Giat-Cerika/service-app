package adminrepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminRepositoryImpl struct {
	db *gorm.DB
}

func NewAdminRepositoryImpl(db *gorm.DB) IAdminRepository {
	return &AdminRepositoryImpl{db: db}
}

// Create implements IAdminRepository.
func (a *AdminRepositoryImpl) Create(ctx context.Context, data *models.User) error {
	return a.db.WithContext(ctx).Create(data).Error
}

// FindUsername implements IAdminRepository.
func (a *AdminRepositoryImpl) FindUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := a.db.WithContext(ctx).Preload("Role").First(&user, "username = ?", username).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// FindRoleAdmin implements IAdminRepository.
func (a *AdminRepositoryImpl) FindRoleAdmin(ctx context.Context) (*models.Role, error) {
	var role models.Role
	if err := a.db.WithContext(ctx).First(&role, "name = ?", "admin").Error; err != nil {
		return nil, err
	}

	return &role, nil
}

// UpdatePhotoAdmin implements IAdminRepository.
func (a *AdminRepositoryImpl) UpdatePhotoAdmin(ctx context.Context, adminID uuid.UUID, photo string) error {
	subQuery := a.db.Select("id").Where("name = ?", "admin").Table("roles")
	return a.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", adminID).
		Where("role_id IN (?)", subQuery).
		Update("photo", photo).Error
}

// FindByAdminID implements IAdminRepository.
func (a *AdminRepositoryImpl) FindByAdminID(ctx context.Context, adminID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := a.db.WithContext(ctx).Preload("Role").First(&user, "id = ?", adminID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// FindSuperAdmin implements IAdminRepository.
func (a *AdminRepositoryImpl) FindAdmin(ctx context.Context, adminId uuid.UUID) (*models.User, error) {
	var admin models.User
	if err := a.db.WithContext(ctx).
		Joins("LEFT JOIN roles ON roles.id = users.role_id").
		Where("users.id = ? AND roles.name = ?", adminId, "admin").
		Preload("Role").
		First(&admin).Error; err != nil {
		return nil, err
	}

	return &admin, nil
}
