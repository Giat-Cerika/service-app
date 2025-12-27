package studentservice

import (
	"context"
	studentrequest "giat-cerika-service/internal/dto/request/student_request"
	"giat-cerika-service/internal/models"
	"mime/multipart"

	"github.com/google/uuid"
)

type IStudentService interface {
	Register(ctx context.Context, req studentrequest.RegisterStudentRequest) error
	Login(ctx context.Context, req studentrequest.LoginStudentRequet) (string, error)
	GetProfile(ctx context.Context, studentId uuid.UUID, token string) (*models.User, error)
	Logout(ctx context.Context, studentID uuid.UUID, token string) error
	CheckTokenBlacklisted(ctx context.Context, token string) (bool, error)
	CheckNisnAndDateOfBirth(ctx context.Context, req studentrequest.CheckNisnAndDateOfBirth) (*models.User, error)
	UpdateNewPasswordStudent(ctx context.Context, studentID uuid.UUID, req studentrequest.UpdatePassword) error
	UpdateProfileStudent(ctx context.Context, studentId uuid.UUID, req studentrequest.UpdateProfileRequest) error
	UpdatePhotoStudent(ctx context.Context, studentId uuid.UUID, photo *multipart.FileHeader) error

	CreateTootBrushStudent(ctx context.Context, studentId uuid.UUID, req studentrequest.CreateTootBrushRequest) error
	GetHitoryToothBrush(ctx context.Context, studentId uuid.UUID, typeTime string, page int, limit int) ([]*models.ToootBrushLog, int, error)
	GetAllStudents(ctx context.Context, page int, limit int, search string) ([]*models.User, int, error)
}
