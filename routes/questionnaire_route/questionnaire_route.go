package questionnaireroute

import (
	questionnairehandler "giat-cerika-service/internal/handlers/questionnaire_handler"
	"giat-cerika-service/internal/middlewares"
	questionnairerepo "giat-cerika-service/internal/repositories/questionnaire_repo"
	questionnaireservice "giat-cerika-service/internal/services/questionnaire_service"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func QuestionnaireRoute(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	questionnaireRepo := questionnairerepo.NewQuestionnaireRepositoryImpl(db)
	questionnaireService := questionnaireservice.NewQuestionnaireServiceImpl(questionnaireRepo, rdb)
	questionnaireHandler := questionnairehandler.NewQuestionnaireHandler(questionnaireService)

	questionnaireGroup := e.Group("", middlewares.JWTMiddleware(rdb), middlewares.RoleMiddleware(strings.ToLower("ADMIN")))
	questionnaireGroup.POST("/create", questionnaireHandler.CreateQuestionnaire)
	questionnaireGroup.GET("/all", questionnaireHandler.GetAllQuestionnaire)
	questionnaireGroup.GET("/:questionnaireId", questionnaireHandler.GetByIdQuestionnaire)
	questionnaireGroup.PUT("/:questionnaireId/edit", questionnaireHandler.UpdateQuestionnaire)
	questionnaireGroup.DELETE("/:questionnaireId/delete", questionnaireHandler.DeleteQuestionnaire)
}
