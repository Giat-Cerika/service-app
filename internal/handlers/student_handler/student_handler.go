package studenthandler

import (
	studentrequest "giat-cerika-service/internal/dto/request/student_request"
	studentresponse "giat-cerika-service/internal/dto/response/student_response"
	studentservice "giat-cerika-service/internal/services/student_service"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	"giat-cerika-service/pkg/constant/response"
	"giat-cerika-service/pkg/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type StudentHandler struct {
	studentService studentservice.IStudentService
}

func NewStudentHandler(service studentservice.IStudentService) *StudentHandler {
	return &StudentHandler{studentService: service}
}

func (s *StudentHandler) RegisterStudent(c echo.Context) error {
	var req studentrequest.RegisterStudentRequest
	req.Name = c.FormValue("name")
	req.Username = c.FormValue("username")
	req.Password = c.FormValue("password")
	req.ConfirmPassword = c.FormValue("confirm_password")
	req.Nisn = c.FormValue("nisn")
	if dateStr := c.FormValue("date_of_birth"); dateStr != "" {
		dateOfBirth, err := time.Parse("02-01-2006", dateStr)
		if err != nil {
			return response.Error(c, http.StatusBadRequest, "invalid date of birth format", err.Error())
		}

		req.DateOfBirth = dateOfBirth
	}

	if v := c.FormValue("age"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			req.Age = &n
		}
	}

	if photo, err := c.FormFile("photo"); err == nil {
		req.Photo = photo
	}

	if v := c.FormValue("class_id"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			req.ClassID = id
		}
	}

	err := s.studentService.Register(c.Request().Context(), req)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to register student")
	}

	return response.Success(c, http.StatusCreated, "Student Registered Successfully", nil)
}

func (s *StudentHandler) LoginStudent(c echo.Context) error {
	var req studentrequest.LoginStudentRequet
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	token, err := s.studentService.Login(c.Request().Context(), req)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to login student")
	}

	return response.Success(c, http.StatusOK, "Login Successfully", map[string]interface{}{
		"access_token": token,
	})
}

func (s *StudentHandler) GetProfileStudent(c echo.Context) error {
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: "+err.Error(), nil)
	}
	studentId := claims.UserID
	authHeader := c.Request().Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	me, err := s.studentService.GetProfile(c.Request().Context(), uuid.MustParse(studentId), token)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get student profile")
	}

	studentResponse := studentresponse.ToStudentResponse(*me)
	return response.Success(c, http.StatusOK, "Get Profile Successfully", studentResponse)
}

func (s *StudentHandler) Logout(c echo.Context) error {
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: "+err.Error(), nil)
	}
	studentId := claims.UserID
	authHeader := c.Request().Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	blacklist, err := s.studentService.CheckTokenBlacklisted(c.Request().Context(), token)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get profile")
	}
	if blacklist {
		return response.Error(c, http.StatusUnauthorized, "unauthorized access", "token blacklisted")
	}
	if token == "" {
		return response.Error(c, http.StatusBadRequest, "bad request: missing token", nil)
	}

	if err := s.studentService.Logout(c.Request().Context(), uuid.MustParse(studentId), token); err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get student profile")
	}

	return response.Success(c, http.StatusOK, "Logout Success", nil)
}

func (s *StudentHandler) CheckNisnAndDateOfBirthStudent(c echo.Context) error {
	var req studentrequest.CheckNisnAndDateOfBirth
	req.Nisn = c.FormValue("nisn")
	if dateStr := c.FormValue("date_of_birth"); dateStr != "" {
		dateOfBirth, err := time.Parse("02-01-2006", dateStr)
		if err != nil {
			return response.Error(c, http.StatusBadRequest, "invalid date of birth format", err.Error())
		}

		req.DateOfBirth = dateOfBirth
	}

	student, err := s.studentService.CheckNisnAndDateOfBirth(c.Request().Context(), req)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to check nisn and date of birth")
	}
	studentResponse := studentresponse.ToStudentResponse(*student)
	return response.Success(c, http.StatusOK, "Check Nisn and Date of Birth Successfully", studentResponse)
}

func (s *StudentHandler) UpdateNewPasswordStudent(c echo.Context) error {
	var req studentrequest.UpdatePassword
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	err := s.studentService.UpdateNewPasswordStudent(c.Request().Context(), req.StudentID, req)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to update password")
	}

	return response.Success(c, http.StatusOK, "Password updated successfully", nil)
}
