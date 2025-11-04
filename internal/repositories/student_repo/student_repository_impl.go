package studentrepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"gorm.io/gorm"
)

type StudentRepositoryImpl struct {
	db *gorm.DB
}

func NewStudentRepositoryImpl(db *gorm.DB) IStudentRepository {
	return &StudentRepositoryImpl{db: db}
}

// Create implements IStudentRepository.
func (s *StudentRepositoryImpl) Create(ctx context.Context, data models.User) error {
	return s.db.WithContext(ctx).Create(data).Error
}
