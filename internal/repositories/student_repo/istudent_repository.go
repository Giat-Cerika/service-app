package studentrepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type IStudentRepository interface {
	Create(ctx context.Context, data *models.User) error
	FindUsernameUnique(ctx context.Context, username string) (string, error)
	FindNisnUnique(ctx context.Context, nisn string) (string, error)
	FindRoleStudent(ctx context.Context) (*models.Role, error)
	UpdatePhotoStudent(ctx context.Context, studentId uuid.UUID, photo string) error

	FindByUsername(ctx context.Context, username string) (*models.User, error)
}
