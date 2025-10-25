package services

import (
	"errors"
	"giat-cerika-service/internal/dto/request"
	"giat-cerika-service/internal/dto/response"
	"giat-cerika-service/internal/repositories"
	"giat-cerika-service/pkg/utils"
)

type adminService struct {
	adminRepo repositories.AdminRepository
}

func NewAdminService(adminRepo repositories.AdminRepository) AdminService {
	return &adminService{adminRepo}
}

func (s *adminService) Login(req request.LoginRequest) (*response.LoginResponse, error) {
	admin, err := s.adminRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(req.Password, admin.Password) {
		return nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(admin.ID, "admin")
	if err != nil {
		return nil, err
	}

	return &response.LoginResponse{
		Token: token,
		User: response.AdminResponse{
			ID:        admin.ID,
			Name:      admin.Name,
			Email:     admin.Email,
			Phone:     admin.Phone,
			Role:      admin.Role,
			CreatedAt: admin.CreatedAt,
		},
	}, nil
}

func (s *adminService) GetProfile(adminID uint) (*response.AdminResponse, error) {
	admin, err := s.adminRepo.FindByID(adminID)
	if err != nil {
		return nil, errors.New("admin not found")
	}

	return &response.AdminResponse{
		ID:        admin.ID,
		Name:      admin.Name,
		Email:     admin.Email,
		Phone:     admin.Phone,
		Role:      admin.Role,
		CreatedAt: admin.CreatedAt,
	}, nil
}

func (s *adminService) UpdateProfile(adminID uint, req request.UpdateAdminRequest) (*response.AdminResponse, error) {
	admin, err := s.adminRepo.FindByID(adminID)
	if err != nil {
		return nil, errors.New("admin not found")
	}

	if req.Name != "" {
		admin.Name = req.Name
	}
	if req.Phone != "" {
		admin.Phone = req.Phone
	}

	if err := s.adminRepo.Update(admin); err != nil {
		return nil, err
	}

	return &response.AdminResponse{
		ID:        admin.ID,
		Name:      admin.Name,
		Email:     admin.Email,
		Phone:     admin.Phone,
		Role:      admin.Role,
		CreatedAt: admin.CreatedAt,
	}, nil
}
