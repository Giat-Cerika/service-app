package utils

import (
	"giat-cerika-service/internal/dto/response"

	"github.com/labstack/echo/v4"
)

func SuccessResponse(c echo.Context, code int, message string, data interface{}) error {
	return c.JSON(code, response.BaseResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c echo.Context, code int, message string) error {
	return c.JSON(code, response.BaseResponse{
		Success: false,
		Message: message,
		Data:    nil,
	})
}

func PaginationResponse(c echo.Context, code int, message string, data interface{}, pagination response.Pagination) error {
	return c.JSON(code, response.PaginationResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}
