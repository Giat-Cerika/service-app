package studentrepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StudentRepositoryImpl struct {
	db *gorm.DB
}

func NewStudentRepositoryImpl(db *gorm.DB) IStudentRepository {
	return &StudentRepositoryImpl{db: db}
}

// Create implements IStudentRepository.
func (s *StudentRepositoryImpl) Create(ctx context.Context, data *models.User) error {
	return s.db.WithContext(ctx).Create(data).Error
}

// FindUsernameUnique implements IStudentRepository.
func (s *StudentRepositoryImpl) FindUsernameUnique(ctx context.Context, username string) (string, error) {
	var student models.User
	if err := s.db.WithContext(ctx).Preload("Role").Preload("Class").First(&student, "username = ?", username).Error; err != nil {
		return "", err
	}

	return student.Username, nil
}

// FindNisnUnique implements IStudentRepository.
func (s *StudentRepositoryImpl) FindNisnUnique(ctx context.Context, nisn string) (string, error) {
	var student models.User
	if err := s.db.WithContext(ctx).Preload("Role").Preload("Class").First(&student, "nisn = ?", nisn).Error; err != nil {
		return "", err
	}

	return *student.Nisn, nil
}

// FindRoleStudent implements IStudentRepository.
func (s *StudentRepositoryImpl) FindRoleStudent(ctx context.Context) (*models.Role, error) {
	var role models.Role
	if err := s.db.WithContext(ctx).First(&role, "name = ?", "student").Error; err != nil {
		return nil, err
	}

	return &role, nil
}

// UpdatePhotoStudent implements IStudentRepository.
func (s *StudentRepositoryImpl) UpdatePhotoStudent(ctx context.Context, studentId uuid.UUID, photo string) error {
	subQuery := s.db.Select("id").Where("name = ?", "student").Table("roles")
	return s.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", studentId).
		Where("role_id IN (?)", subQuery).
		Update("photo", photo).Error
}

// FindByUsername implements IStudentRepository.
func (s *StudentRepositoryImpl) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var student models.User
	if err := s.db.WithContext(ctx).Preload("Role").Preload("Class").First(&student, "username = ?", username).Error; err != nil {
		return nil, err
	}

	return &student, nil
}
