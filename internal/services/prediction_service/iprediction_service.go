package predictionservice

import (
	"context"
	predictionrequest "giat-cerika-service/internal/dto/request/prediction_request"
	"giat-cerika-service/internal/models"
)

type IPredictionService interface {
	Create(ctx context.Context, req predictionrequest.CreatePredictionRequest) error
	GetAllPrediction(ctx context.Context, page, limit int, search string) ([]*models.Prediction, int, error)
}
