package studentrequest

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type RegisterStudentRequest struct {
	Name            string                `form:"name" json:"name"`
	Username        string                `form:"username" json:"username"`
	Password        string                `form:"password" json:"password"`
	ConfirmPassword string                `form:"confirm_password" json:"confirm_password"`
	Nisn            string                `form:"nisn" json:"nisn"`
	DateOfBirth     time.Time             `form:"date_of_birth" json:"date_of_birth"`
	Age             *int                  `form:"age" json:"age"`
	Photo           *multipart.FileHeader `form:"photo" json:"photo"`
	ClassID         uuid.UUID             `form:"class_id" json:"class_id"`
}

type LoginStudentRequet struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

type CheckNisnAndDateOfBirth struct {
	Nisn        string    `form:"nisn" json:"nisn"`
	DateOfBirth time.Time `form:"date_of_birth" json:"date_of_birth"`
}

type UpdatePassword struct {
	StudentID       uuid.UUID `form:"student_id" json:"student_id"`
	NewPassword     string    `form:"new_password" json:"new_password"`
	ConfirmPassword string    `form:"confirm_password" json:"confirm_password"`
}
