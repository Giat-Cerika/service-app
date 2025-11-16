package videoroute

import (
	videohandler "giat-cerika-service/internal/handlers/video_handler"
	"giat-cerika-service/internal/middlewares"
	videorepo "giat-cerika-service/internal/repositories/video_repo"
	videoservice "giat-cerika-service/internal/services/video_service"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func VideoRoute(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	videoRepo := videorepo.NewVideoRepositoryImpl(db)
	videoService := videoservice.NewVideoServiceImpl(videoRepo, rdb)
	videoHandler := videohandler.NewVideoHandler(videoService)

	videoGroup := e.Group("", middlewares.JWTMiddleware(rdb), middlewares.RoleMiddleware(strings.ToLower("ADMIN")))
	videoGroup.POST("/create", videoHandler.CreateVideo)
	videoGroup.GET("/all", videoHandler.GetAllVideo)
	videoGroup.GET("/:videoId", videoHandler.GetByIdVideo)
	videoGroup.PUT("/:videoId/edit", videoHandler.UpdateVideo)
	videoGroup.DELETE("/:videoId/delete", videoHandler.DeleteVideo)
}
