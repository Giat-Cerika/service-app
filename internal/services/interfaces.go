package services

import (
	"giat-cerika-service/internal/dto/request"
	"giat-cerika-service/internal/dto/response"
)

type AdminService interface {
	Login(req request.LoginRequest) (*response.LoginResponse, error)
	GetProfile(adminID uint) (*response.AdminResponse, error)
	UpdateProfile(adminID uint, req request.UpdateAdminRequest) (*response.AdminResponse, error)
}

type StudentService interface {
	Register(req request.RegisterStudentRequest) (*response.StudentResponse, error)
	Login(req request.LoginRequest) (*response.LoginResponse, error)
	GetProfile(studentID uint) (*response.StudentResponse, error)
	UpdateProfile(studentID uint, req request.UpdateStudentRequest) (*response.StudentResponse, error)
}
