package classrepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type IClassRepository interface {
	Create(ctx context.Context, data *models.Class) error
	FindByNameClass(ctx context.Context, NameClass string) (*models.Class, error)
	FindAll(ctx context.Context, limit, offset int, search string) ([]*models.Class, int, error)
	FindById(ctx context.Context, classId uuid.UUID) (*models.Class, error)
	Update(ctx context.Context, classId uuid.UUID, data *models.Class) error
	Delete(ctx context.Context, classId uuid.UUID) error

	GetAllPublic(ctx context.Context) ([]*models.Class, error)
}
