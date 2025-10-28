package adminservice

import (
	"context"
	"errors"
	"fmt"
	datasources "giat-cerika-service/internal/dataSources"
	adminrequest "giat-cerika-service/internal/dto/request/admin_request"
	"giat-cerika-service/internal/models"
	adminrepo "giat-cerika-service/internal/repositories/admin_repo"
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

type AdminServiceImpl struct {
	adminRepo adminrepo.IAdminRepository
	rdb       *redis.Client
	cld       datasources.CloudinaryService
}

func NewAdminServiceImpl(adminRepo adminrepo.IAdminRepository, rdb *redis.Client, cld datasources.CloudinaryService) IAdminService {
	return &AdminServiceImpl{adminRepo: adminRepo, rdb: rdb, cld: cld}
}

func fileAdminToBytes(fh *multipart.FileHeader) ([]byte, error) {
	f, err := fh.Open()
	if err != nil {
		return nil, err
	}

	defer f.Close()
	return io.ReadAll(f)
}

// Register implements IAdminService.
func (a *AdminServiceImpl) Register(ctx context.Context, req adminrequest.RegisterAdminRequest) error {
	existsUsername, err := a.adminRepo.FindUsername(ctx, req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Failed to get user username", 500)
	}

	if strings.TrimSpace(req.Username) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Username is Required", 400)
	}
	if strings.TrimSpace(req.Password) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Password is Required", 400)
	}
	if req.Photo == nil {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Photo is Required", 400)
	}

	if existsUsername != nil {
		return errorresponse.NewCustomError(errorresponse.ErrExists, "Username Already Exists", 409)
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Failed to hashing password", 400)
	}
	role, err := a.adminRepo.FindRoleAdmin(ctx)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Failed to get role admin", 500)
	}

	admin := &models.User{
		ID:       uuid.New(),
		Username: req.Username,
		Password: hashed,
		RoleID:   role.ID,
		Role:     *role,
		Status:   1,
	}

	if err := a.adminRepo.Create(ctx, admin); err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to create admin", 500)
	}

	if req.Photo != nil {
		if binner, err := fileAdminToBytes(req.Photo); err == nil && len(binner) > 0 {
			pay := payload.ImageUploadPayload{
				ID:        admin.ID,
				Type:      "single",
				FileBytes: binner,
				Folder:    "giat_ceria/photo_admin",
				Filename:  fmt.Sprintf("admin_%s_photo", admin.ID.String()),
			}

			_ = rabbitmq.PublishToQueue("", rabbitmq.SendImageProfileAdminQueueName, pay)
		}
	}

	return nil
}
