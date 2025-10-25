package handlers

import (
	"giat-cerika-service/internal/dto/request"
	"giat-cerika-service/internal/services"
	"giat-cerika-service/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type StudentHandler struct {
	studentService services.StudentService
}

func NewStudentHandler(studentService services.StudentService) *StudentHandler {
	return &StudentHandler{studentService}
}

func (h *StudentHandler) Register(c echo.Context) error {
	var req request.RegisterStudentRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	result, err := h.studentService.Register(req)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusCreated, "Registration successful", result)
}

func (h *StudentHandler) Login(c echo.Context) error {
	var req request.LoginRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	result, err := h.studentService.Login(req)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "Login successful", result)
}

func (h *StudentHandler) GetProfile(c echo.Context) error {
	studentID := utils.GetUserIDFromToken(c)

	result, err := h.studentService.GetProfile(studentID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "Profile retrieved successfully", result)
}

func (h *StudentHandler) UpdateProfile(c echo.Context) error {
	studentID := utils.GetUserIDFromToken(c)

	var req request.UpdateStudentRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request")
	}

	result, err := h.studentService.UpdateProfile(studentID, req)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "Profile updated successfully", result)
}

func (h *StudentHandler) Logout(c echo.Context) error {
	return utils.SuccessResponse(c, http.StatusOK, "Logout successful", nil)
}
