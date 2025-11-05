package classhandler

import (
	classrequest "giat-cerika-service/internal/dto/request/class_request"
	classresponse "giat-cerika-service/internal/dto/response/class_response"
	classservice "giat-cerika-service/internal/services/class_service"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	"giat-cerika-service/pkg/constant/response"
	"giat-cerika-service/pkg/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ClassHandler struct {
	classService classservice.IClassService
}

func NewClassHandler(service classservice.IClassService) *ClassHandler {
	return &ClassHandler{classService: service}
}

func (ch *ClassHandler) CreateClass(c echo.Context) error {
	var req classrequest.CreateClassRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	err := ch.classService.CreateClass(c.Request().Context(), req)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to create class")
	}

	return response.Success(c, http.StatusOK, "Class Created Succssfully", nil)
}

func (ch *ClassHandler) GetAllClass(c echo.Context) error {
	pageInt, limitInt := utils.ParsePaginationParams(c, 10)
	search := c.QueryParam("search")

	classes, total, err := ch.classService.GetAllClass(c.Request().Context(), pageInt, limitInt, search)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get classes")
	}

	meta := utils.BuildPaginationMeta(c, pageInt, limitInt, total)
	data := make([]classresponse.ClassResponse, len(classes))
	for i, class := range classes {
		data[i] = classresponse.ToClassResponse(*class)
	}

	return response.PaginatedSuccess(c, http.StatusOK, "Get All Classes Successfully", data, meta)
}

func (ch *ClassHandler) GetByIdClass(c echo.Context) error {
	classId, err := uuid.Parse(c.Param("classId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	class, err := ch.classService.GetByIdClass(c.Request().Context(), classId)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get class")
	}

	res := classresponse.ToClassResponse(*class)

	return response.Success(c, http.StatusOK, "Get Class Successfully", res)
}

func (ch *ClassHandler) UpdateClass(c echo.Context) error {
	classId, err := uuid.Parse(c.Param("classId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	var req classrequest.UpdateClassRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	err = ch.classService.UpdateClass(c.Request().Context(), classId, req)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to update class")
	}

	return response.Success(c, http.StatusOK, "Class Updated Successfully", nil)
}

func (ch *ClassHandler) DeleteClass(c echo.Context) error {
	classId, err := uuid.Parse(c.Param("classId"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	if err := ch.classService.DeleteClass(c.Request().Context(), classId); err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to delete class")
	}

	return response.Success(c, http.StatusOK, "Class Deleted Successfully", nil)
}
