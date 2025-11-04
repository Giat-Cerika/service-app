package studentrequest

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type RegisterStudentRequest struct {
	Name            string                `form:"name" json:"name"`
	Username        string                `form:"username" json:"username"`
	Password        string                `form:"password" json:"password"`
	ConfirmPassword string                `form:"confirm_password" json:"confirm_password"`
	Nisn            string                `form:"nisn" json:"nisn"`
	DateOfBirth     string                `form:"date_of_birth" json:"date_of_birth"`
	Age             int                   `form:"age" json:"age"`
	Photo           *multipart.FileHeader `form:"photo" json:"photo"`
	ClassID         uuid.UUID             `form:"class_id" json:"class_id"`
}
