package detailquizroute

import (
	detailquizhandler "giat-cerika-service/internal/handlers/detail_quiz_handler"
	"giat-cerika-service/internal/middlewares"
	detailquizrepo "giat-cerika-service/internal/repositories/detail_quiz_repo"
	detailquizservice "giat-cerika-service/internal/services/detail_quiz_service"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func MaterialRoute(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	detail_quizRepo := detailquizrepo.NewMaterialRepositoryImpl(db)
	detail_quizService := detailquizservice.NewMaterialServiceImpl(detail_quizRepo, rdb)
	detail_quizHandler := detailquizhandler.NewMaterialHandler(detail_quizService)

	detail_quizGroup := e.Group("", middlewares.JWTMiddleware(rdb), middlewares.RoleMiddleware(strings.ToLower("ADMIN")))
	detail_quizGroup.POST("/create", detail_quizHandler.CreateMaterial)
	detail_quizGroup.GET("/all", detail_quizHandler.GetAllMaterial)
	detail_quizGroup.GET("/:detail_quizId", detail_quizHandler.GetByIdMaterial)
	detail_quizGroup.PUT("/:detail_quizId/edit", detail_quizHandler.UpdateMaterial)
	detail_quizGroup.DELETE("/:detail_quizId/delete", detail_quizHandler.DeleteMaterial)
}
