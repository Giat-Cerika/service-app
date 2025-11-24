package videorepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type IVideoRepository interface {
	Create(ctx context.Context, data *models.Video) error
	FindByTitle(ctx context.Context, NameVideo string) (*models.Video, error)
	FindAll(ctx context.Context, limit, offset int, search string) ([]*models.Video, int, error)
	FindById(ctx context.Context, videoId uuid.UUID) (*models.Video, error)
	Update(ctx context.Context, videoId uuid.UUID, data *models.Video) error
	Delete(ctx context.Context, videoId uuid.UUID) error

	FindAllLatest(ctx context.Context) ([]*models.Video, error)
	FindAllPublic(ctx context.Context, limit, offset int, search string) ([]*models.Video, int, error)
}
