package response

import (
	"giat-cerika-service/internal/models"
	"time"

	"github.com/google/uuid"
)

type AdminResponse struct {
	ID        uuid.UUID    `json:"id"`
	Name      string       `json:"name"`
	Username  string       `json:"username"`
	Phone     string       `json:"phone"`
	Role      RoleResponse `json:"role"`
	CreatedAt time.Time    `json:"created_at"`
}

func ToAdminResponse(admin models.User) AdminResponse {
	dataRole := ToRoleResponse(admin.Role)

	return AdminResponse{ID: admin.ID, Name: admin.Name, Username: admin.Username, Role: dataRole}
}
