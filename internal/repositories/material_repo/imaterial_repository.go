package materialrepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type IMaterialRepository interface {
	Create(ctx context.Context, data *models.Materials) error
	FindByTitle(ctx context.Context, NameMaterial string) (*models.Materials, error)
	FindAll(ctx context.Context, limit, offset int, search string) ([]*models.Materials, int, error)
	FindById(ctx context.Context, materialId uuid.UUID) (*models.Materials, error)
	Update(ctx context.Context, materialId uuid.UUID, data *models.Materials) error
	Delete(ctx context.Context, materialId uuid.UUID) error

	UpdateCoverMateri(ctx context.Context, materiId uuid.UUID, cover string) error
	CreateImage(ctx context.Context, image *models.Image) error
	CreateGallery(ctx context.Context, materiId uuid.UUID, imageId uuid.UUID, alt string) error
	DeleteGalleryByMateriId(ctx context.Context, materiId uuid.UUID) error

	FindAllLatest(ctx context.Context) ([]*models.Materials, error)
	FindAllPublic(ctx context.Context, limit, offset int, search string) ([]*models.Materials, int, error)
}
