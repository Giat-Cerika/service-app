package adminresponse

import (
	"giat-cerika-service/internal/models"
	"giat-cerika-service/pkg/utils"

	"github.com/google/uuid"
)

type AdminResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Status    int       `json:"status"`
	Photo     string    `json:"photo"`
	Role      string    `json:"role"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func ToAdminResponse(admin models.User) AdminResponse {
	return AdminResponse{
		ID:        admin.ID,
		Username:  admin.Username,
		Status:    admin.Status,
		Photo:     admin.Photo,
		Role:      admin.Role.Name,
		CreatedAt: utils.FormatDate(admin.CreatedAt),
		UpdatedAt: utils.FormatDate(admin.UpdatedAt),
	}
}
