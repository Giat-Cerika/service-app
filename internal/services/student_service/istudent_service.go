package studentservice

import (
	"context"
	studentrequest "giat-cerika-service/internal/dto/request/student_request"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type IStudentService interface {
	Register(ctx context.Context, req studentrequest.RegisterStudentRequest) error
	Login(ctx context.Context, req studentrequest.LoginStudentRequet) (string, error)
	GetProfile(ctx context.Context, studentId uuid.UUID, token string) (*models.User, error)
	Logout(ctx context.Context, studentID uuid.UUID, token string) error
	CheckTokenBlacklisted(ctx context.Context, token string) (bool, error)
	CheckNisnAndDateOfBirth(ctx context.Context, req studentrequest.CheckNisnAndDateOfBirth) (*models.User, error)
}
