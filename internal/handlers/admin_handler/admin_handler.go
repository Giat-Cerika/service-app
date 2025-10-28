package adminhandler

import (
	adminrequest "giat-cerika-service/internal/dto/request/admin_request"
	adminservice "giat-cerika-service/internal/services/admin_service"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	"giat-cerika-service/pkg/constant/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	adminService adminservice.IAdminService
}

func NewAdminHandler(service adminservice.IAdminService) *AdminHandler {
	return &AdminHandler{adminService: service}
}

func (a *AdminHandler) RegisterAdmin(c echo.Context) error {
	var req adminrequest.RegisterAdminRequest
	req.Username = c.FormValue("username")
	req.Password = c.FormValue("password")
	if photo, err := c.FormFile("photo"); err == nil {
		req.Photo = photo
	}

	err := a.adminService.Register(c.Request().Context(), req)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to create admin")
	}

	return response.Success(c, http.StatusCreated, "Admin Created Successfully", nil)

}
