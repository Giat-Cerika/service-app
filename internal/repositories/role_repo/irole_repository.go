package rolerepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type IRoleRepository interface {
	Create(ctx context.Context, data *models.Role) error
	FindByName(ctx context.Context, name string) (*models.Role, error)

	FindAll(ctx context.Context, limit, offset int, search string) ([]*models.Role, int, error)
	FindById(ctx context.Context, roleId uuid.UUID) (*models.Role, error)
	Update(ctx context.Context, roleId uuid.UUID, data *models.Role) error
	Delete(ctx context.Context, roleId uuid.UUID) error
}
