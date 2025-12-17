package quizsessionroute

import (
	quizsessionhandler "giat-cerika-service/internal/handlers/quiz_session_handler"
	"giat-cerika-service/internal/middlewares"
	quizrepo "giat-cerika-service/internal/repositories/quiz_repo"
	quizsessionrepo "giat-cerika-service/internal/repositories/quiz_session_repo"
	studentrepo "giat-cerika-service/internal/repositories/student_repo"
	quizsessionservice "giat-cerika-service/internal/services/quiz_session_service"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func QuizSessionRoute(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	qsRepo := quizsessionrepo.NewQuizSessionRepositoryImpl(db)
	quizRepo := quizrepo.NewQuizRepositoryImpl(db)
	studentRepo := studentrepo.NewStudentRepositoryImpl(db)
	qsService := quizsessionservice.NewQuizSessionServiceImpl(qsRepo, quizRepo, studentRepo, rdb)
	qsHandler := quizsessionhandler.NewQuizSessionHandler(qsService)

	qsGroup := e.Group("", middlewares.JWTMiddleware(rdb), middlewares.RoleMiddleware(strings.ToLower("STUDENT")))
	qsGroup.POST("/assign-code-quiz/:quizId", qsHandler.AssignCodeQuiz)
	qsGroup.PUT("/start-quiz/:quizSessionId", qsHandler.StartedQuiz)
	qsGroup.GET("/quiz-duration/:quizSessionId", qsHandler.GetDuration)
	qsGroup.POST("/quiz-submit/:quizSessionId", qsHandler.SubmitQuizSession)
	qsGroup.GET("/quiz-question/:quizSessionId", qsHandler.GetQuizQuestionByOrderMode)

	qsAdmin := e.Group("", middlewares.JWTMiddleware(rdb), middlewares.RoleMiddleware(strings.ToLower("ADMIN")))
	qsAdmin.GET("/all-student", qsHandler.GetQuizSessionStudent)
}
