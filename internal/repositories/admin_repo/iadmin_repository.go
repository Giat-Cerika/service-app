package adminrepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type IAdminRepository interface {
	Create(ctx context.Context, data *models.User) error
	FindUsername(ctx context.Context, username string) (*models.User, error)
	FindRoleAdmin(ctx context.Context) (*models.Role, error)
	UpdatePhotoAdmin(ctx context.Context, adminID uuid.UUID, photo string) error
}
