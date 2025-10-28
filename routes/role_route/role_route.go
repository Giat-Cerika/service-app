package roleroute

import (
	rolehandler "giat-cerika-service/internal/handlers/role_handler"
	rolerepo "giat-cerika-service/internal/repositories/role_repo"
	roleservice "giat-cerika-service/internal/services/role_service"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RoleRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	roleRepo := rolerepo.NewRoleRepositoryImpl(db)
	roleService := roleservice.NewRoleServiceImpl(roleRepo, rdb)
	roleHandler := rolehandler.NewRoleHandler(roleService)

	e.POST("/create", roleHandler.CreateRole)
	e.GET("/all", roleHandler.GetAllRole)
	e.GET("/:roleId", roleHandler.GetByIdRole)
	e.PUT("/:roleId/edit", roleHandler.UpdateRole)
	e.DELETE("/:roleId/delete", roleHandler.DeleteRole)
}
