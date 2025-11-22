package routes

import (
	datasources "giat-cerika-service/internal/dataSources"
	adminroute "giat-cerika-service/routes/admin_route"
	classroute "giat-cerika-service/routes/class_route"
	roleroute "giat-cerika-service/routes/role_route"
	studentroute "giat-cerika-service/routes/student_route"
	videoroute "giat-cerika-service/routes/video_route"
	materialroute "giat-cerika-service/routes/material_route"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Routes(e *echo.Echo, db *gorm.DB, rdb *redis.Client, cldSvc *datasources.CloudinaryService) {
	v1 := e.Group("/api/v1")
	roleroute.RoleRoutes(v1.Group("/role"), db, rdb)
	adminroute.AdminRoutes(v1.Group("/admin"), db, rdb, cldSvc)
	classroute.ClassRoutes(v1.Group("/class"), db, rdb)
	studentroute.StudentRoutes(v1.Group("/student"), db, rdb, cldSvc)
	videoroute.VideoRoutes(v1.Group("/video"), db, rdb)
	materialroute.MaterialRoute(v1.Group("/material"), db, rdb, cldSvc)
}
