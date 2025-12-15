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

func QuizRoute(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	quizRepo := quizrepo.NewQuizRepositoryImpl(db)
	qtRepo := quizrepo.NewQuizTypeRepositoryImpl(db)
	quizService := quizservice.NewQuizServiceImpl(quizRepo, qtRepo, rdb)
	quizHandler := quizhandler.NewQuizHandler(quizService)

	quizGroup := e.Group("", middlewares.JWTMiddleware(rdb), middlewares.RoleMiddleware(strings.ToLower("ADMIN")))
	quizGroup.POST("/create", quizHandler.CreateQuiz)
	quizGroup.GET("/all", quizHandler.GetQuizAll)
	quizGroup.GET("/:quizId", quizHandler.GetQuizByID)
	quizGroup.PUT("/:quizId/edit", quizHandler.UpdateQuiz)
	quizGroup.DELETE("/:quizId/delete", quizHandler.DeleteQuiz)
	quizGroup.PUT("/:quizId/update-status", quizHandler.UpdateStatusQuiz)
	quizGroup.PUT("/:quizId/update-question-order-mode", quizHandler.UpdateQuestionOrderMode)
}
