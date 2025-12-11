package quizroute

import (
	quizhandler "giat-cerika-service/internal/handlers/quiz_handler"
	"giat-cerika-service/internal/middlewares"
	quizrepo "giat-cerika-service/internal/repositories/quiz_repo"
	quizservice "giat-cerika-service/internal/services/quiz_service"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func QuizTypeRoute(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	qtRepo := quizrepo.NewQuizTypeRepositoryImpl(db)
	qtService := quizservice.NewQuizTypeServiceImpl(qtRepo, rdb)
	qtHandler := quizhandler.NewQuizTypeHandler(qtService)

	qtGroup := e.Group("", middlewares.JWTMiddleware(rdb), middlewares.RoleMiddleware(strings.ToLower("ADMIN")))
	qtGroup.POST("/create", qtHandler.CreateQuizType)
	qtGroup.GET("/all", qtHandler.GetAllQuizType)
	qtGroup.GET("/:quizTypeid", qtHandler.GetQuizTypeByID)
	qtGroup.PUT("/:quizTypeid/edit", qtHandler.UpdateQuizType)
	qtGroup.DELETE("/:quizTypeid/delete", qtHandler.DeleteQuizType)
}
