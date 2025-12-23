package models

import (
	"time"

	"github.com/google/uuid"
)

type Prediction struct {
	ID                 uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PatientName        string           `gorm:"type:varchar(255)" json:"patient_name"`
	Age                int              `gorm:"type:int" json:"age"`
	DateOfEvaluation   time.Time        `gorm:"type:timestamptz" json:"date_of_evaluation"`
	Confidence         string           `gorm:"type:varchar(100)" json:"confidence"`
	ConfidenceDetailID uuid.UUID        `gorm:"type:uuid" json:"-"`
	ConfidenceDetail   ConfidenceDetail `gorm:"foreignKey:ConfidenceDetailID;constraint:OnDelete:CASCADE;" json:"confidence_detail"`
	CariesRiskID       uuid.UUID        `gorm:"type:uuid" json:"-"`
	CariesRisk         CariesRisk       `gorm:"foreignKey:CariesRiskID;constraint:OnDelete:CASCADE;" json:"caries_risk"`
	Description        string           `gorm:"type:text" json:"description"`
	Result             string           `gorm:"type:varchar(100)" json:"result"`
	Score              int              `gorm:"type:int" json:"score"`
	CreatedAt          time.Time        `gorm:"autoCreateTime"`
	UpdatedAt          time.Time        `gorm:"autoUpdateTime"`
}

type ConfidenceDetail struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	Low    int `gorm:"type:int" json:"low"`
	Medium int `gorm:"type:int" json:"medium"`
	High   int `gorm:"type:int" json:"high"`
}

type CariesRisk struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	AttitudeAndStatus int `gorm:"type:int" json:"attitude_and_status"`
	CariesHistory     int `gorm:"type:int" json:"caries_history"`
	Fluoride          int `gorm:"type:int" json:"fluoride"`
	ModifyingFactor   int `gorm:"type:int" json:"modifying_factor"`

	DietID   uuid.UUID  `gorm:"type:uuid"`
	PlaqueID *uuid.UUID `gorm:"type:uuid"`
	SalivaID *uuid.UUID `gorm:"type:uuid"`

	Diet   DietDetail    `gorm:"foreignKey:DietID;constraint:OnDelete:CASCADE;" json:"diet"`
	Plaque *PlaqueOption `gorm:"foreignKey:PlaqueID;constraint:OnDelete:CASCADE;" json:"plaque,omitempty"`
	Saliva *SalivaOption `gorm:"foreignKey:SalivaID;constraint:OnDelete:CASCADE;" json:"saliva,omitempty"`
}

type DietDetail struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	Acid  int `gorm:"type:int" json:"acid"`
	Sugar int `gorm:"type:int" json:"sugar"`
}

type PlaqueOption struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	Maturity *int `gorm:"type:int" json:"maturity"`
	Ph       *int `gorm:"type:int" json:"ph"`
}

type SalivaOption struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	RestingSalivaID    *uuid.UUID `gorm:"type:uuid;nullable" json:"-"`
	StimulatedSalivaID *uuid.UUID `gorm:"type:uuid;nullable" json:"-"`

	RestingSaliva    *RestingSaliva    `gorm:"foreignKey:RestingSalivaID;constraint:OnDelete:CASCADE;" json:"resting_saliva,omitempty"`
	StimulatedSaliva *StimulatedSaliva `gorm:"foreignKey:StimulatedSalivaID;constraint:OnDelete:CASCADE;" json:"stimulated_saliva,omitempty"`
}

type RestingSaliva struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	Hydration int `gorm:"type:int" json:"hydration"`
	Viscosity int `gorm:"type:int" json:"viscosity"`
	Ph        int `gorm:"type:int" json:"ph"`
}

type StimulatedSaliva struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	Quantity  int `gorm:"type:int" json:"quantity"`
	Ph        int `gorm:"type:int" json:"ph"`
	Buffering int `gorm:"type:int" json:"buffering"`
}
