package classroute

import (
	classhandler "giat-cerika-service/internal/handlers/class_handler"
	"giat-cerika-service/internal/middlewares"
	classrepo "giat-cerika-service/internal/repositories/class_repo"
	classservice "giat-cerika-service/internal/services/class_service"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func ClassRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	classRepo := classrepo.NewClassRepositoryImpl(db)
	classService := classservice.NewClassServiceImpl(classRepo, rdb)
	classHandler := classhandler.NewClassHandler(classService)

	classGroup := e.Group("", middlewares.JWTMiddleware(rdb), middlewares.RoleMiddleware(strings.ToLower("ADMIN")))
	classGroup.POST("/create", classHandler.CreateClass)
	classGroup.GET("/all", classHandler.GetAllClass)
	classGroup.GET("/:classId", classHandler.GetByIdClass)
	classGroup.PUT("/:classId/edit", classHandler.UpdateClass)
	classGroup.DELETE("/:classId/delete", classHandler.DeleteClass)
}
