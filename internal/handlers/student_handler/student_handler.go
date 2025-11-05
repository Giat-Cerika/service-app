package studenthandler

import (
	studentrequest "giat-cerika-service/internal/dto/request/student_request"
	studentservice "giat-cerika-service/internal/services/student_service"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	"giat-cerika-service/pkg/constant/response"
	"net/http"
	"strconv"
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
