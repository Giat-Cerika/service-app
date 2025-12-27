package quizhistoryroute

import (
	quizhistoryhandler "giat-cerika-service/internal/handlers/quiz_history_handler"
	"giat-cerika-service/internal/middlewares"
	quizhistoryrepo "giat-cerika-service/internal/repositories/quiz_history_repo"
	quizrepo "giat-cerika-service/internal/repositories/quiz_repo"
	studentrepo "giat-cerika-service/internal/repositories/student_repo"
	quizhistoryservice "giat-cerika-service/internal/services/quiz_history_service"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func QuizHistoryRoute(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	quizHistoryRepo := quizhistoryrepo.NewQuizHistoryRepositoryImpl(db)
	studentRepo := studentrepo.NewStudentRepositoryImpl(db)
	quizRepo := quizrepo.NewQuizRepositoryImpl(db)
	quizHistoryService := quizhistoryservice.NewQuizHistoryServiceImpl(quizHistoryRepo, studentRepo, quizRepo, rdb)
	quizHistoryHandler := quizhistoryhandler.NewQuizHistoryHandler(quizHistoryService)

	qhGruop := e.Group("", middlewares.JWTMiddleware(rdb), middlewares.RoleMiddleware(strings.ToLower("STUDENT")))
	qhGruop.GET("/my-history", quizHistoryHandler.GetHistoryQuizStudent)
	qhGruop.GET("/question-history/:quizHistoryId", quizHistoryHandler.GetAllQuestionHistory)

	qhAdmin := e.Group("", middlewares.JWTMiddleware(rdb), middlewares.RoleMiddleware(strings.ToLower("ADMIN")))
	qhAdmin.GET("/all-student-history", quizHistoryHandler.GetHistoryQuizByQuizID)

}
