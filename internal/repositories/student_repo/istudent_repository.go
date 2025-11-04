package studentrepo

import (
	"context"
	"giat-cerika-service/internal/models"
)

type IStudentRepository interface {
	Create(ctx context.Context, data models.User) error
}
