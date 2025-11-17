package classservice

import (
	"context"
	classrequest "giat-cerika-service/internal/dto/request/class_request"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type IClassService interface {
	CreateClass(ctx context.Context, req classrequest.CreateClassRequest) error
	GetAllClass(ctx context.Context, page, limit int, search string) ([]*models.Class, int, error)
	GetByIdClass(ctx context.Context, classId uuid.UUID) (*models.Class, error)
	UpdateClass(ctx context.Context, classId uuid.UUID, req classrequest.UpdateClassRequest) error
	DeleteClass(ctx context.Context, classId uuid.UUID) error

	GetAllPublic(ctx context.Context) ([]*models.Class, error)
}
