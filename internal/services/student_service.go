package services

import (
	"errors"
	"giat-cerika-service/internal/dto/request"
	"giat-cerika-service/internal/dto/response"
	"giat-cerika-service/internal/models"
	"giat-cerika-service/internal/repositories"
	"giat-cerika-service/pkg/utils"
)

type studentService struct {
	studentRepo repositories.StudentRepository
}

func NewStudentService(studentRepo repositories.StudentRepository) StudentService {
	return &studentService{studentRepo}
}

func (s *studentService) Register(req request.RegisterStudentRequest) (*response.StudentResponse, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	student := &models.Student{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Phone:    req.Phone,
		NIM:      req.NIM,
		Major:    req.Major,
	}

	if err := s.studentRepo.Create(student); err != nil {
		return nil, errors.New("failed to register student")
	}

	return &response.StudentResponse{
		ID:        student.ID,
		Name:      student.Name,
		Email:     student.Email,
		Phone:     student.Phone,
		NIM:       student.NIM,
		Major:     student.Major,
		CreatedAt: student.CreatedAt,
	}, nil
}

func (s *studentService) Login(req request.LoginRequest) (*response.LoginResponse, error) {
	student, err := s.studentRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(req.Password, student.Password) {
		return nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(student.ID, "student")
	if err != nil {
		return nil, err
	}

	return &response.LoginResponse{
		Token: token,
		User: response.StudentResponse{
			ID:        student.ID,
			Name:      student.Name,
			Email:     student.Email,
			Phone:     student.Phone,
			NIM:       student.NIM,
			Major:     student.Major,
			CreatedAt: student.CreatedAt,
		},
	}, nil
}

func (s *studentService) GetProfile(studentID uint) (*response.StudentResponse, error) {
	student, err := s.studentRepo.FindByID(studentID)
	if err != nil {
		return nil, errors.New("student not found")
	}

	return &response.StudentResponse{
		ID:        student.ID,
		Name:      student.Name,
		Email:     student.Email,
		Phone:     student.Phone,
		NIM:       student.NIM,
		Major:     student.Major,
		CreatedAt: student.CreatedAt,
	}, nil
}

func (s *studentService) UpdateProfile(studentID uint, req request.UpdateStudentRequest) (*response.StudentResponse, error) {
	student, err := s.studentRepo.FindByID(studentID)
	if err != nil {
		return nil, errors.New("student not found")
	}

	if req.Name != "" {
		student.Name = req.Name
	}
	if req.Phone != "" {
		student.Phone = req.Phone
	}
	if req.Major != "" {
		student.Major = req.Major
	}

	if err := s.studentRepo.Update(student); err != nil {
		return nil, err
	}

	return &response.StudentResponse{
		ID:        student.ID,
		Name:      student.Name,
		Email:     student.Email,
		Phone:     student.Phone,
		NIM:       student.NIM,
		Major:     student.Major,
		CreatedAt: student.CreatedAt,
	}, nil
}
