package adminroute

import (
	datasources "giat-cerika-service/internal/dataSources"
	adminhandler "giat-cerika-service/internal/handlers/admin_handler"
	"giat-cerika-service/internal/middlewares"
	adminrepo "giat-cerika-service/internal/repositories/admin_repo"
	adminservice "giat-cerika-service/internal/services/admin_service"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func AdminRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client, cld *datasources.CloudinaryService) {
	adminRepo := adminrepo.NewAdminRepositoryImpl(db)
	adminService := adminservice.NewAdminServiceImpl(adminRepo, rdb, *cld)
	adminHandler := adminhandler.NewAdminHandler(adminService)

	publicGroup := e.Group("")
	publicGroup.POST("/register", adminHandler.RegisterAdmin)
	publicGroup.POST("/login", adminHandler.LoginAdmin)

	protectedGroup := e.Group("", middlewares.JWTMiddleware(rdb), middlewares.RoleMiddleware("admin"))
	protectedGroup.GET("/me", adminHandler.GetProfileAdmin)
	protectedGroup.POST("/logout", adminHandler.LogoutAdmin)
}
