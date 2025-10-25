package handlers

import (
	"giat-cerika-service/internal/dto/request"
	"giat-cerika-service/internal/services"
	"giat-cerika-service/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	adminService services.AdminService
}

func NewAdminHandler(adminService services.AdminService) *AdminHandler {
	return &AdminHandler{adminService}
}

func (h *AdminHandler) Login(c echo.Context) error {
	var req request.LoginRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	result, err := h.adminService.Login(req)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "Login successful", result)
}

func (h *AdminHandler) GetProfile(c echo.Context) error {
	adminID := utils.GetUserIDFromToken(c)

	result, err := h.adminService.GetProfile(adminID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "Profile retrieved successfully", result)
}

func (h *AdminHandler) UpdateProfile(c echo.Context) error {
	adminID := utils.GetUserIDFromToken(c)

	var req request.UpdateAdminRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request")
	}

	result, err := h.adminService.UpdateProfile(adminID, req)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "Profile updated successfully", result)
}

func (h *AdminHandler) Logout(c echo.Context) error {
	return utils.SuccessResponse(c, http.StatusOK, "Logout successful", nil)
}
