package adminservice

import (
	"context"
	adminrequest "giat-cerika-service/internal/dto/request/admin_request"
)

type IAdminService interface {
	Register(ctx context.Context, req adminrequest.RegisterAdminRequest) error
}
