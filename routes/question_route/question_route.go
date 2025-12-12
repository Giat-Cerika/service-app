package questionroute

import (
	datasources "giat-cerika-service/internal/dataSources"
	questionhandler "giat-cerika-service/internal/handlers/question_handler"
	"giat-cerika-service/internal/middlewares"
	answerrepo "giat-cerika-service/internal/repositories/answer_repo"
	questionrepo "giat-cerika-service/internal/repositories/question_repo"
	quizrepo "giat-cerika-service/internal/repositories/quiz_repo"
	questionservice "giat-cerika-service/internal/services/question_service"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func QuestionRoute(e *echo.Group, db *gorm.DB, rdb *redis.Client, cld *datasources.CloudinaryService) {
	questionRepo := questionrepo.NewQuestionRepositoryImpl(db)
	quizRepo := quizrepo.NewQuizRepositoryImpl(db)
	answerRepo := answerrepo.NewAnswerRepositoryImpl(db)
	questionService := questionservice.NewQuestionServiceImpl(questionRepo, quizRepo, answerRepo, rdb, *cld)
	questionHandler := questionhandler.NewQuestionHandler(questionService)

	questionGroup := e.Group("", middlewares.JWTMiddleware(rdb), middlewares.RoleMiddleware(strings.ToLower("ADMIN")))
	questionGroup.POST("/create", questionHandler.CreateQuestion)
	questionGroup.GET("/all", questionHandler.GetAllQuestion)
	questionGroup.GET("/:questionId", questionHandler.GetByIdQuestion)
	questionGroup.PUT("/:questionId/edit", questionHandler.UpdateQuestion)
	questionGroup.DELETE("/:questionId/delete", questionHandler.DeleteQuestion)
}
