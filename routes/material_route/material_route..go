package materialroute

import (
	datasources "giat-cerika-service/internal/dataSources"
	materialhandler "giat-cerika-service/internal/handlers/material_handler"
	"giat-cerika-service/internal/middlewares"
	adminrepo "giat-cerika-service/internal/repositories/admin_repo"
	materialrepo "giat-cerika-service/internal/repositories/material_repo"
	materialservice "giat-cerika-service/internal/services/material_service"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func MaterialRoute(e *echo.Group, db *gorm.DB, rdb *redis.Client, cld *datasources.CloudinaryService) {
	materialRepo := materialrepo.NewMaterialRepositoryImpl(db)
	adminRepo := adminrepo.NewAdminRepositoryImpl(db)
	materialService := materialservice.NewMaterialServiceImpl(materialRepo, adminRepo, rdb, *cld)
	materialHandler := materialhandler.NewMaterialHandler(materialService)

	materialGroup := e.Group("", middlewares.JWTMiddleware(rdb), middlewares.RoleMiddleware(strings.ToLower("ADMIN")))
	materialGroup.POST("/create", materialHandler.CreateMaterial)
	materialGroup.GET("/all", materialHandler.GetAllMaterial)
	materialGroup.GET("/:materialId", materialHandler.GetByIdMaterial)
	materialGroup.PUT("/:materialId/edit", materialHandler.UpdateMaterial)
	materialGroup.DELETE("/:materialId/delete", materialHandler.DeleteMaterial)
}
