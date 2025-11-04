package adminservice

import (
	"context"
	adminrequest "giat-cerika-service/internal/dto/request/admin_request"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type IAdminService interface {
	Register(ctx context.Context, req adminrequest.RegisterAdminRequest) error
	Login(ctx context.Context, req adminrequest.LoginAdminRequest) (string, error)
	GetProfile(ctx context.Context, adminId uuid.UUID, token string) (*models.User, error)
	Logout(ctx context.Context, adminID uuid.UUID, token string) error
	CheckTokenBlacklisted(ctx context.Context, token string) (bool, error)
}
