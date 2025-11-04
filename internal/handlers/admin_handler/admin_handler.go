package adminhandler

import (
	adminrequest "giat-cerika-service/internal/dto/request/admin_request"
	adminresponse "giat-cerika-service/internal/dto/response/admin_response"
	adminservice "giat-cerika-service/internal/services/admin_service"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	"giat-cerika-service/pkg/constant/response"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func (a *AdminHandler) LoginAdmin(c echo.Context) error {
	var req adminrequest.LoginAdminRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	token, err := a.adminService.Login(c.Request().Context(), req)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "invalid login admin")
	}

	return response.Success(c, http.StatusOK, "Login Successfully", map[string]interface{}{
		"access_token": token,
	})

}

func (a *AdminHandler) GetProfileAdmin(c echo.Context) error {
	adminToken := c.Get("user")
	if adminToken == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: token invalid or expired", nil)
	}

	admin, ok := adminToken.(*jwt.Token)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: token invalid or expired", nil)
	}

	claims := admin.Claims.(jwt.MapClaims)
	adminID := claims["user_id"].(string)
	authHeader := c.Request().Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	me, err := a.adminService.GetProfile(c.Request().Context(), uuid.MustParse(adminID), token)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "invalid to get profile")
	}

	adminResponse := adminresponse.ToAdminResponse(*me)
	return response.Success(c, http.StatusOK, "Get Profile Successfully", adminResponse)
}

func (a *AdminHandler) LogoutAdmin(c echo.Context) error {
	adminToken := c.Get("user")
	if adminToken == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: token invalid or expired", nil)
	}

	admin, ok := adminToken.(*jwt.Token)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: token invalid or expired", nil)
	}

	claims := admin.Claims.(jwt.MapClaims)
	adminID := claims["user_id"].(string)
	authHeader := c.Request().Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	blackList, err := a.adminService.CheckTokenBlacklisted(c.Request().Context(), token)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "invalid to get blacklisted")
	}
	if blackList {
		return response.Error(c, http.StatusUnauthorized, "Your'e logged out", nil)
	}

	if token != "" {
		if err := a.adminService.Logout(c.Request().Context(), uuid.MustParse(adminID), token); err != nil {
			if customErr, ok := errorresponse.AsCustomErr(err); ok {
				return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
			}
			return response.Error(c, http.StatusInternalServerError, err.Error(), "invalid to blacklist token")
		}
	} else {
		return response.Error(c, http.StatusUnauthorized, "Token Is empty", nil)
	}

	return response.Success(c, http.StatusOK, "Logout Successfully", nil)
}
