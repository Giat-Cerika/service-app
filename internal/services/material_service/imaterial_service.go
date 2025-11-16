package materialservice

import (
	"context"
	materialrequest "giat-cerika-service/internal/dto/request/material_request"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type IMaterialService interface {
	CreateMaterial(ctx context.Context, req materialrequest.CreateMaterialRequest) error
	GetAllMaterial(ctx context.Context, page, limit int, search string) ([]*models.Materials, int, error)
	GetByIdMaterial(ctx context.Context, materialId uuid.UUID) (*models.Materials, error)
	UpdateMaterial(ctx context.Context, materialId uuid.UUID, req materialrequest.UpdateMaterialRequest) error
	DeleteMaterial(ctx context.Context, materialId uuid.UUID) error
}
