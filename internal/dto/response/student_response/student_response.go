package studentresponse

import (
	"giat-cerika-service/internal/models"
	"giat-cerika-service/pkg/utils"

	"github.com/google/uuid"
)

type StudentResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Username    string    `json:"username"`
	Nisn        string    `json:"nisn"`
	DateOfBirth string    `json:"date_of_birth"`
	Age         int       `json:"age"`
	Photo       string    `json:"photo"`
	Role        string    `json:"role"`
	Class       string    `json:"class"`
	Status      int       `json:"status"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

type AllStudentResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Username    string    `json:"username"`
	Nisn        string    `json:"nisn"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

func ToStudentResponse(student models.User) StudentResponse {
	return StudentResponse{
		ID:          student.ID,
		Name:        *student.Name,
		Username:    student.Username,
		Nisn:        *student.Nisn,
		DateOfBirth: utils.FormatDate(*student.DateOfBirth),
		Age:         student.Age,
		Photo:       student.Photo,
		Role:        student.Role.Name,
		Class:       student.Class.NameClass,
		Status:      student.Status,
		CreatedAt:   utils.FormatDate(student.CreatedAt),
		UpdatedAt:   utils.FormatDate(student.UpdatedAt),
	}
}

func ToAllStudentResponse(student models.User) AllStudentResponse {
	var name, nisn string

	if student.Name != nil {
		name = *student.Name
	}

	if student.Nisn != nil {
		nisn = *student.Nisn
	}

	return AllStudentResponse{
		ID:        student.ID,
		Name:      name,
		Username:  student.Username,
		Nisn:      nisn,
		CreatedAt: utils.FormatDate(student.CreatedAt),
		UpdatedAt: utils.FormatDate(student.UpdatedAt),
	}
}
