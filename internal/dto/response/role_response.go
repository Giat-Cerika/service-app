package response

import (
	"giat-cerika-service/internal/models"
)

type RoleResponse struct {
	Name     string `gorm:"type:varchar(100);not null" json:"name"`
	Email    string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
	Phone    string `gorm:"type:varchar(20)" json:"phone"`
	NIM      string `gorm:"type:varchar(50);uniqueIndex" json:"nim"`
	Major    string `gorm:"type:varchar(100)" json:"major"`
}

func ToRoleResponse(role models.Role) RoleResponse {
	return RoleResponse{Name: role.Name, Major: role.Major}
}