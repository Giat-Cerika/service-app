package routes

import (
	datasources "giat-cerika-service/internal/dataSources"
	roleroute "giat-cerika-service/routes/role_route"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Routes(e *echo.Echo, db *gorm.DB, rdb *redis.Client, cldSvc *datasources.CloudinaryService) {
	v1 := e.Group("/api/v1")
	roleroute.RoleRoutes(v1.Group("/role"), db, rdb)
}
