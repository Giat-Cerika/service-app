package predictionrepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
)

type IPredictionRepository interface {
	CreatePrediction(ctx context.Context, data *models.Prediction) error
	CreateConfidenceDetail(ctx context.Context, data *models.ConfidenceDetail) error
	CreateDietDetail(ctx context.Context, data *models.DietDetail) error
	CreatePlaqueOption(ctx context.Context, data *models.PlaqueOption) error
	CreateCariesRisk(ctx context.Context, data *models.CariesRisk) error

	CreateRestingSaliva(ctx context.Context, data *models.RestingSaliva) error
	CreateStimulatedSaliva(ctx context.Context, data *models.StimulatedSaliva) error
	CreateSalivaOption(ctx context.Context, data *models.SalivaOption) error

	GetAllPrediction(ctx context.Context, limit, offset int, search string) ([]*models.Prediction, int, error)
	GetByIdPrediction(ctx context.Context, predictionId uuid.UUID) (*models.Prediction, error)
	DeletePrediction(ctx context.Context, predictionId uuid.UUID) error
}
