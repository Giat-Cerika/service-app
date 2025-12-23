package predictionroute

import (
	predictionhandler "giat-cerika-service/internal/handlers/prediction_handler"
	"giat-cerika-service/internal/middlewares"
	predictionrepo "giat-cerika-service/internal/repositories/prediction_repo"
	predictionservice "giat-cerika-service/internal/services/prediction_service"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func PredictionRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client) {
	predictRepo := predictionrepo.NewPredictionRepositoryImpl(db)
	predicService := predictionservice.NewPredictionServiceImpl(predictRepo, rdb)
	predictHandler := predictionhandler.NewPredictionHandler(predicService)

	predictGroup := e.Group("", middlewares.JWTMiddleware(rdb), middlewares.RoleMiddleware(strings.ToLower("ADMIN")))
	predictGroup.POST("/save", predictHandler.CreatePrediction)
	predictGroup.GET("/all", predictHandler.GetAllPredictions)
}
