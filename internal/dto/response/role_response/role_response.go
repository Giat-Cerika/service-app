package roleresponse

import (
	"giat-cerika-service/internal/models"
	"giat-cerika-service/pkg/utils"

	"github.com/google/uuid"
)

type RoleResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func ToRoleResponse(role models.Role) RoleResponse {
	return RoleResponse{
		ID:        role.ID,
		Name:      role.Name,
		CreatedAt: utils.FormatDate(role.CreatedAt),
		UpdatedAt: utils.FormatDate(role.UpdatedAt),
	}
}
