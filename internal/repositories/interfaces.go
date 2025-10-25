package repositories

import (
	"giat-cerika-service/internal/models"
)

type AdminRepository interface {
	FindByEmail(email string) (*models.Admin, error)
	FindByID(id uint) (*models.Admin, error)
	Update(admin *models.Admin) error
}

type StudentRepository interface {
	Create(student *models.Student) error
	FindByEmail(email string) (*models.Student, error)
	FindByID(id uint) (*models.Student, error)
	Update(student *models.Student) error
}
