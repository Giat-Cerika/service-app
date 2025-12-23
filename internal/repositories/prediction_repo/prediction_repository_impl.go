package predictionrepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"gorm.io/gorm"
)

type PredictionRepositoryImpl struct {
	db *gorm.DB
}

func NewPredictionRepositoryImpl(db *gorm.DB) IPredictionRepository {
	return &PredictionRepositoryImpl{db: db}
}

func (p *PredictionRepositoryImpl) DB() *gorm.DB {
	return p.db
}

// CreateCariesRisk implements [IPredictionRepository].
func (p *PredictionRepositoryImpl) CreateCariesRisk(ctx context.Context, data *models.CariesRisk) error {
	return p.db.WithContext(ctx).Create(data).Error
}

// CreateConfidenceDetail implements [IPredictionRepository].
func (p *PredictionRepositoryImpl) CreateConfidenceDetail(ctx context.Context, data *models.ConfidenceDetail) error {
	return p.db.WithContext(ctx).Create(data).Error
}

// CreateDietDetail implements [IPredictionRepository].
func (p *PredictionRepositoryImpl) CreateDietDetail(ctx context.Context, data *models.DietDetail) error {
	return p.db.WithContext(ctx).Create(data).Error
}

// CreatePlaqueOption implements [IPredictionRepository].
func (p *PredictionRepositoryImpl) CreatePlaqueOption(ctx context.Context, data *models.PlaqueOption) error {
	return p.db.WithContext(ctx).Create(data).Error
}

// CreateRestingSaliva implements [IPredictionRepository].
func (p *PredictionRepositoryImpl) CreateRestingSaliva(ctx context.Context, data *models.RestingSaliva) error {
	return p.db.WithContext(ctx).Create(data).Error
}

// CreateSalivaOption implements [IPredictionRepository].
func (p *PredictionRepositoryImpl) CreateSalivaOption(ctx context.Context, data *models.SalivaOption) error {
	return p.db.WithContext(ctx).Create(data).Error
}

// CreateStimulatedSaliva implements [IPredictionRepository].
func (p *PredictionRepositoryImpl) CreateStimulatedSaliva(ctx context.Context, data *models.StimulatedSaliva) error {
	return p.db.WithContext(ctx).Create(data).Error
}

// Create implements [IPredictionRepository].
func (p *PredictionRepositoryImpl) CreatePrediction(ctx context.Context, data *models.Prediction) error {
	return p.db.WithContext(ctx).Create(data).Error
}

// GetAllPrediction implements [IPredictionRepository].
func (p *PredictionRepositoryImpl) GetAllPrediction(ctx context.Context, limit int, offset int, search string) ([]*models.Prediction, int, error) {
	var (
		predictions []*models.Prediction
		count       int64
	)

	query := p.db.WithContext(ctx).
		Model(&models.Prediction{}).
		Preload("ConfidenceDetail").
		Preload("CariesRisk").
		Preload("CariesRisk.Diet").
		Preload("CariesRisk.Plaque").
		Preload("CariesRisk.Saliva").
		Preload("CariesRisk.Saliva.RestingSaliva").
		Preload("CariesRisk.Saliva.StimulatedSaliva")
	if search != "" {
		query = query.Where("patient_name ILIKE ?", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order("date_of_evaluation DESC").Find(&predictions).Error; err != nil {
		return nil, 0, err
	}

	return predictions, int(count), nil
}
