package predictionresponse

import (
	"giat-cerika-service/internal/models"
	"giat-cerika-service/pkg/utils"

	"github.com/google/uuid"
)

type ConfidenceDetailResponse struct {
	Low    int `json:"low"`
	Medium int `json:"medium"`
	High   int `json:"high"`
}

type CariesRiskResponse struct {
	AttitudeAndStatus int             `json:"attitude_and_status"`
	CariesHistory     int             `json:"caries_history"`
	Fluoride          int             `json:"fluoride"`
	ModifyingFactor   int             `json:"modifying_factor"`
	Diet              DietResponse    `json:"diet"`
	Plaque            *PlaqueResponse `json:"plaque,omitempty"`
	Saliva            *SalivaResponse `json:"saliva,omitempty"`
}

type DietResponse struct {
	Acid  int `json:"acid"`
	Sugar int `json:"sugar"`
}

type PlaqueResponse struct {
	Maturity *int `json:"maturity,omitempty"`
	Ph       *int `json:"ph,omitempty"`
}

type SalivaResponse struct {
	RestingSaliva    *RestingSalivaResponse    `json:"resting_saliva,omitempty"`
	StimulatedSaliva *StimulatedSalivaResponse `json:"stimulated_saliva,omitempty"`
}

type RestingSalivaResponse struct {
	Hydration int `json:"hydration"`
	Viscosity int `json:"viscosity"`
	Ph        int `json:"ph"`
}

type StimulatedSalivaResponse struct {
	Quantity  int `json:"quantity"`
	Ph        int `json:"ph"`
	Buffering int `json:"buffering"`
}

type PredictionResponse struct {
	ID               uuid.UUID                `json:"id"`
	PatientName      string                   `json:"patient_name"`
	Age              int                      `json:"age"`
	DateOfEvaluation string                   `json:"date_of_evaluation"`
	Confidence       string                   `json:"confidence"`
	ConfidenceDetail ConfidenceDetailResponse `json:"confidence_detail"`
	Score            int                      `json:"score"`
	Result           string                   `json:"result"`
	Description      string                   `json:"description"`
	CariesRisk       CariesRiskResponse       `json:"caries_risk"`
}

func ToPredictionResponse(p models.Prediction) PredictionResponse {

	resp := PredictionResponse{
		ID:               p.ID,
		PatientName:      p.PatientName,
		Age:              p.Age,
		DateOfEvaluation: utils.FormatDate(p.DateOfEvaluation),
		Confidence:       p.Confidence,
		Score:            p.Score,
		Result:           p.Result,
		Description:      p.Description,
		ConfidenceDetail: ConfidenceDetailResponse{
			Low:    p.ConfidenceDetail.Low,
			Medium: p.ConfidenceDetail.Medium,
			High:   p.ConfidenceDetail.High,
		},
		CariesRisk: CariesRiskResponse{
			AttitudeAndStatus: p.CariesRisk.AttitudeAndStatus,
			CariesHistory:     p.CariesRisk.CariesHistory,
			Fluoride:          p.CariesRisk.Fluoride,
			ModifyingFactor:   p.CariesRisk.ModifyingFactor,
			Diet: DietResponse{
				Acid:  p.CariesRisk.Diet.Acid,
				Sugar: p.CariesRisk.Diet.Sugar,
			},
		},
	}

	// OPTIONAL PLAQUE
	if p.CariesRisk.Plaque != nil {
		resp.CariesRisk.Plaque = &PlaqueResponse{
			Maturity: p.CariesRisk.Plaque.Maturity,
			Ph:       p.CariesRisk.Plaque.Ph,
		}
	}

	// OPTIONAL SALIVA
	if p.CariesRisk.Saliva != nil {
		saliva := &SalivaResponse{}

		if p.CariesRisk.Saliva.RestingSaliva != nil {
			saliva.RestingSaliva = &RestingSalivaResponse{
				Hydration: p.CariesRisk.Saliva.RestingSaliva.Hydration,
				Viscosity: p.CariesRisk.Saliva.RestingSaliva.Viscosity,
				Ph:        p.CariesRisk.Saliva.RestingSaliva.Ph,
			}
		}

		if p.CariesRisk.Saliva.StimulatedSaliva != nil {
			saliva.StimulatedSaliva = &StimulatedSalivaResponse{
				Quantity:  p.CariesRisk.Saliva.StimulatedSaliva.Quantity,
				Ph:        p.CariesRisk.Saliva.StimulatedSaliva.Ph,
				Buffering: p.CariesRisk.Saliva.StimulatedSaliva.Buffering,
			}
		}

		resp.CariesRisk.Saliva = saliva
	}

	return resp
}

type PredictionByStudentResponse struct {
	ID               uuid.UUID `json:"id"`
	PatientName      string    `json:"patient_name"`
	Age              int       `json:"age"`
	DateOfEvaluation string    `json:"date_of_evaluation"`
	Confidence       string    `json:"confidence"`
	Result           string    `json:"result"`
	Score            int       `json:"score"`
	Description      string    `json:"description"`
	Suggestion       string    `json:"suggestion"`
	CreatedAt        string    `json:"created_at"`
}

func ToPredictionByStudentResponse(ps models.PredictHistory) PredictionByStudentResponse {
	return PredictionByStudentResponse{
		ID:               ps.ID,
		PatientName:      ps.Prediction.PatientName,
		Age:              ps.Prediction.Age,
		DateOfEvaluation: utils.FormatDate(ps.Prediction.DateOfEvaluation),
		Confidence:       ps.Prediction.Confidence,
		Result:           ps.Prediction.Result,
		Score:            ps.Prediction.Score,
		Description:      ps.Prediction.Description,
		Suggestion:       ps.Suggestion,
		CreatedAt:        utils.FormatDate(ps.CreatedAt),
	}
}
