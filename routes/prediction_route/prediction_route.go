package predictionroute

import (
	predictionhandler "giat-cerika-service/internal/handlers/prediction_handler"
	"giat-cerika-service/internal/middlewares"
	predictionrepo "giat-cerika-service/internal/repositories/prediction_repo"
	studentrepo "giat-cerika-service/internal/repositories/student_repo"
	predictionservice "giat-cerika-service/internal/services/prediction_service"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func PredictionRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	predictRepo := predictionrepo.NewPredictionRepositoryImpl(db)
	studentRepo := studentrepo.NewStudentRepositoryImpl(db)
	predicService := predictionservice.NewPredictionServiceImpl(predictRepo, studentRepo, rdb)
	predictHandler := predictionhandler.NewPredictionHandler(predicService)

	predictGroup := e.Group("", middlewares.JWTMiddleware(rdb), middlewares.RoleMiddleware(strings.ToLower("ADMIN")))
	predictGroup.POST("/save", predictHandler.CreatePrediction)
	predictGroup.GET("/all", predictHandler.GetAllPredictions)
	predictGroup.DELETE("/:predictionId/delete", predictHandler.DeletePrediction)

	predictGroup.POST("/send-prediction", predictHandler.SendPredictToStudent)

	predictStudent := e.Group("", middlewares.JWTMiddleware(rdb), middlewares.RoleMiddleware(strings.ToLower("STUDENT")))
	predictStudent.GET("/my-prediction", predictHandler.GetPredictByStudent)
}
