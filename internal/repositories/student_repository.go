package repositories

import (
	"giat-cerika-service/internal/models"

	"gorm.io/gorm"
)

type studentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) StudentRepository {
	return &studentRepository{db}
}

func (r *studentRepository) Create(student *models.Student) error {
	return r.db.Create(student).Error
}

func (r *studentRepository) FindByEmail(email string) (*models.Student, error) {
	var student models.Student
	if err := r.db.Where("email = ?", email).First(&student).Error; err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) FindByID(id uint) (*models.Student, error) {
	var student models.Student
	if err := r.db.First(&student, id).Error; err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) Update(student *models.Student) error {
	return r.db.Save(student).Error
}
