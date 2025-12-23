package predictionservice

import (
	"context"
	"encoding/json"
	"fmt"
	"giat-cerika-service/configs"
	predictionrequest "giat-cerika-service/internal/dto/request/prediction_request"
	"giat-cerika-service/internal/models"
	predictionrepo "giat-cerika-service/internal/repositories/prediction_repo"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type PredictionServiceImpl struct {
	predictionRepo predictionrepo.IPredictionRepository
	rdb            *redis.Client
}

func NewPredictionServiceImpl(predicRepo predictionrepo.IPredictionRepository, rdb *redis.Client) IPredictionService {
	return &PredictionServiceImpl{predictionRepo: predicRepo, rdb: rdb}
}

func (p *PredictionServiceImpl) invalidateCachePrediction(ctx context.Context) {
	iter := p.rdb.Scan(ctx, 0, "predictions:*", 0).Iterator()
	for iter.Next(ctx) {
		p.rdb.Del(ctx, iter.Val())
	}
}

// Create implements [IPredictionService].
func (p *PredictionServiceImpl) Create(ctx context.Context, req predictionrequest.CreatePredictionRequest) error {
	repoImpl := p.predictionRepo.(*predictionrepo.PredictionRepositoryImpl)
	db := repoImpl.DB()

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		confidence := models.ConfidenceDetail{
			ID:     uuid.New(),
			Low:    req.ConfidenceDetail.Low,
			Medium: req.ConfidenceDetail.Medium,
			High:   req.ConfidenceDetail.High,
		}
		if err := tx.Create(&confidence).Error; err != nil {
			return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to save data confidence", 500)
		}

		diet := models.DietDetail{
			ID:    uuid.New(),
			Acid:  req.CariesRisk.Diet.Acid,
			Sugar: req.CariesRisk.Diet.Sugar,
		}
		if err := tx.Create(&diet).Error; err != nil {
			return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to save data diet", 500)
		}

		var plaqueID *uuid.UUID
		if req.CariesRisk.Plaque != nil {
			plaque := models.PlaqueOption{
				ID:       uuid.New(),
				Maturity: req.CariesRisk.Plaque.Maturity,
				Ph:       req.CariesRisk.Plaque.Ph,
			}
			if err := tx.Create(&plaque).Error; err != nil {
				return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to save data plaque", 500)
			}
			plaqueID = &plaque.ID
		}

		var salivaID *uuid.UUID

		if req.CariesRisk.Saliva != nil {
			var restingID *uuid.UUID
			var stimulatedID *uuid.UUID

			if req.CariesRisk.Saliva.RestingSaliva != nil {
				resting := models.RestingSaliva{
					ID:        uuid.New(),
					Hydration: req.CariesRisk.Saliva.RestingSaliva.Hydration,
					Viscosity: req.CariesRisk.Saliva.RestingSaliva.Viscosity,
					Ph:        req.CariesRisk.Saliva.RestingSaliva.Ph,
				}
				if err := tx.Create(&resting).Error; err != nil {
					return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to save data resting saliva", 500)
				}
				restingID = &resting.ID
			}

			if req.CariesRisk.Saliva.StimulatedSaliva != nil {
				stimulated := models.StimulatedSaliva{
					ID:        uuid.New(),
					Quantity:  req.CariesRisk.Saliva.StimulatedSaliva.Quantity,
					Ph:        req.CariesRisk.Saliva.StimulatedSaliva.Ph,
					Buffering: req.CariesRisk.Saliva.StimulatedSaliva.Buffering,
				}
				if err := tx.Create(&stimulated).Error; err != nil {
					return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to save data stimulated saliva", 500)
				}
				stimulatedID = &stimulated.ID
			}

			saliva := models.SalivaOption{
				ID:                 uuid.New(),
				RestingSalivaID:    restingID,
				StimulatedSalivaID: stimulatedID,
			}
			if err := tx.Create(&saliva).Error; err != nil {
				return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to save data saliva", 500)
			}
			salivaID = &saliva.ID
		}

		caries := models.CariesRisk{
			ID:                uuid.New(),
			AttitudeAndStatus: req.CariesRisk.AttitudeAndStatus,
			CariesHistory:     req.CariesRisk.CariesHistory,
			Fluoride:          req.CariesRisk.Fluoride,
			ModifyingFactor:   req.CariesRisk.ModifyingFactor,
			DietID:            diet.ID,
			PlaqueID:          plaqueID,
			SalivaID:          salivaID,
		}
		if err := tx.Create(&caries).Error; err != nil {
			return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to save data caries risk", 500)
		}

		// ===============================
		// 6. Prediction (FINAL)
		// ===============================
		locJakarta, _ := time.LoadLocation("Asia/Jakarta")
		nowJakarta := time.Now().In(locJakarta)
		prediction := models.Prediction{
			ID:                 uuid.New(),
			PatientName:        req.PatientName,
			Age:                req.Age,
			DateOfEvaluation:   nowJakarta,
			Confidence:         req.Confidence,
			Score:              req.Score,
			Result:             req.Result,
			Description:        req.Description,
			ConfidenceDetailID: confidence.ID,
			CariesRiskID:       caries.ID,
		}

		if err := tx.Create(&prediction).Error; err != nil {
			return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to save data prediction", 500)
		}

		p.invalidateCachePrediction(ctx)

		return nil
	})
}

// GetAllPrediction implements [IPredictionService].
func (p *PredictionServiceImpl) GetAllPrediction(ctx context.Context, page int, limit int, search string) ([]*models.Prediction, int, error) {
	cacheKey := fmt.Sprintf("predictions:search:%s:page:%d:limit:%d", search, page, limit)
	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var result struct {
			Data  []*models.Prediction `json:"data"`
			Total int                  `json:"total"`
		}
		if json.Unmarshal([]byte(cached), &result) == nil {
			return result.Data, result.Total, nil
		}
	}

	offset := (page - 1) * limit

	items, total, err := p.predictionRepo.GetAllPrediction(ctx, limit, offset, search)
	if err != nil {
		return nil, 0, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get predictions", 500)
	}

	if len(items) == 0 {
		items = []*models.Prediction{}
	}

	buf, _ := json.Marshal(map[string]any{
		"data":  items,
		"total": total,
	})

	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)

	return items, total, nil
}
