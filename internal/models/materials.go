package models

import (
	"time"

	"github.com/google/uuid"
)

type Materials struct {
	ID             uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title          string           `gorm:"type:varchar(255);index"`
	Description    string           `gorm:"type:text"`
	Cover          string           `gorm:"type:varchar(255)"`
	MaterialImages []MaterialImages `gorm:"foreignKey:MaterialID;constraint:OnDelete:CASCADE"`
	CreatedBy      uuid.UUID        `gorm:"type:uuid"`
	User           User             `gorm:"foreignKey:CreatedBy"`
	CreatedAt      time.Time        `gorm:"autoCreateTime"`
	UpdatedAt      time.Time        `gorm:"autoUpdateTime"`
}

type MaterialImages struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	MaterialID uuid.UUID `gorm:"type:uuid;index" json:"material_id"`
	Material   Materials `gorm:"foreignKey:MaterialID;OnDelete:CASCADE" json:"material"`
	ImageID    uuid.UUID `gorm:"type:uuid;index" json:"image_id"`
	Image      Image     `gorm:"foreignKey:ImageID" json:"image"`
	AltText    string    `gorm:"type:varchar(255)" json:"alt_text"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
