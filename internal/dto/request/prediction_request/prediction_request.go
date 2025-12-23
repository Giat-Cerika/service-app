package predictionrequest

import "time"

type CreatePredictionRequest struct {
	PatientName      string    `json:"patient_name"`
	Age              int       `json:"age"`
	DateOfEvaluation time.Time `json:"date_of_evaluation"`
	Confidence       string    `json:"confidence"`
	Score            int       `json:"score"`
	Result           string    `json:"result"`
	Description      string    `json:"description"`

	ConfidenceDetail ConfidenceDetailRequest `json:"confidence_detail"`
	CariesRisk       CariesRiskRequest       `json:"caries_risk"`
}

type ConfidenceDetailRequest struct {
	Low    int `json:"low"`
	Medium int `json:"medium"`
	High   int `json:"high"`
}

type CariesRiskRequest struct {
	AttitudeAndStatus int `json:"attitude_and_status"`
	CariesHistory     int `json:"caries_history"`
	Fluoride          int `json:"fluoride"`
	ModifyingFactor   int `json:"modifying_factor"`

	Diet   DietRequest    `json:"diet"`
	Plaque *PlaqueRequest `json:"plaque,omitempty"`
	Saliva *SalivaRequest `json:"saliva,omitempty"`
}

type DietRequest struct {
	Acid  int `json:"acid"`
	Sugar int `json:"sugar"`
}

type PlaqueRequest struct {
	Maturity *int `json:"maturity"`
	Ph       *int `json:"ph"`
}

type SalivaRequest struct {
	RestingSaliva    *RestingSalivaRequest    `json:"resting_saliva,omitempty"`
	StimulatedSaliva *StimulatedSalivaRequest `json:"stimulated_saliva,omitempty"`
}

type RestingSalivaRequest struct {
	Hydration int `json:"hydration"`
	Viscosity int `json:"viscosity"`
	Ph        int `json:"ph"`
}

type StimulatedSalivaRequest struct {
	Quantity  int `json:"quantity"`
	Ph        int `json:"ph"`
	Buffering int `json:"buffering"`
}
