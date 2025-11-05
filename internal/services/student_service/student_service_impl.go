package studentservice

import (
	"context"
	"errors"
	"fmt"
	datasources "giat-cerika-service/internal/dataSources"
	studentrequest "giat-cerika-service/internal/dto/request/student_request"
	"giat-cerika-service/internal/models"
	classrepo "giat-cerika-service/internal/repositories/class_repo"
	studentrepo "giat-cerika-service/internal/repositories/student_repo"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	rabbitmq "giat-cerika-service/pkg/constant/rabbitMq"
	"giat-cerika-service/pkg/utils"
	"giat-cerika-service/pkg/workers/payload"
	"io"
	"mime/multipart"
	"strings"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type StudentServiceImpl struct {
	studenRepo studentrepo.IStudentRepository
	classRepo  classrepo.IClassRepository
	rdb        *redis.Client
	cld        datasources.CloudinaryService
}

func NewStudentServiceImpl(studentRepo studentrepo.IStudentRepository, classRepo classrepo.IClassRepository, rdb *redis.Client, cld datasources.CloudinaryService) IStudentService {
	return &StudentServiceImpl{studenRepo: studentRepo, classRepo: classRepo, rdb: rdb, cld: cld}
}

func fileStudentToBytes(fh *multipart.FileHeader) ([]byte, error) {
	f, err := fh.Open()
	if err != nil {
		return nil, err
	}

	defer f.Close()
	return io.ReadAll(f)
}

// Register implements IStudentService.
func (s *StudentServiceImpl) Register(ctx context.Context, req studentrequest.RegisterStudentRequest) error {
	uniqueUsername, err := s.studenRepo.FindUsernameUnique(ctx, req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get unique username", 500)
	}
	uniqueNisn, err := s.studenRepo.FindNisnUnique(ctx, req.Nisn)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get unique nisn", 500)
	}

	if strings.TrimSpace(req.Name) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Name is required", 400)
	}
	if strings.TrimSpace(req.Username) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Username is required", 400)
	}
	if strings.TrimSpace(req.Password) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Password is required", 400)
	}
	if strings.TrimSpace(req.ConfirmPassword) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Confirm Password is required", 400)
	}
	if req.Password != req.ConfirmPassword {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Password and Confirm Password doesn't match", 409)
	}
	if strings.TrimSpace(req.Nisn) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Nisn is required", 400)
	}
	if req.DateOfBirth.IsZero() {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Date Of Birth is required", 400)
	}
	if req.Age == nil {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Age is required", 400)
	}
	if req.Photo == nil {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Photo is required", 400)
	}
	if req.ClassID == uuid.Nil {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "class is required", 400)
	}
	if uniqueUsername != "" {
		return errorresponse.NewCustomError(errorresponse.ErrExists, "Username already exists", 409)
	}
	if uniqueNisn != "" {
		return errorresponse.NewCustomError(errorresponse.ErrExists, "Nisn already exists", 409)
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "failed to hashing password", 400)
	}

	role, err := s.studenRepo.FindRoleStudent(ctx)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get role", 500)
	}

	class, err := s.classRepo.FindById(ctx, req.ClassID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "class not found", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get class", 500)
	}

	newStudent := &models.User{
		ID:          uuid.New(),
		Name:        &req.Name,
		Username:    req.Username,
		Password:    hashed,
		Nisn:        &req.Nisn,
		DateOfBirth: &req.DateOfBirth,
		Age:         *req.Age,
		RoleID:      role.ID,
		ClassID:     &class.ID,
		Status:      1,
	}

	if err := s.studenRepo.Create(ctx, newStudent); err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to create student", 500)
	}

	if req.Photo != nil {
		if binner, err := fileStudentToBytes(req.Photo); err == nil && len(binner) > 0 {
			pay := payload.ImageUploadPayload{
				ID:        newStudent.ID,
				Type:      "single",
				FileBytes: binner,
				Folder:    "giat_ceria/photo_student",
				Filename:  fmt.Sprintf("student_%s_photo", newStudent.ID.String()),
			}

			_ = rabbitmq.PublishToQueue("", rabbitmq.SendImageProfileStudentQueueName, pay)
		}
	}

	return nil
}
