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
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	MaterialID uuid.UUID `gorm:"type:uuid;index"`
	Material   Materials `gorm:"foreignKey:MaterialID;constraint:OnDelete:CASCADE"`
	ImageID    uuid.UUID `gorm:"type:uuid;index"`
	Image      Image     `gorm:"foreignKey:ImageID"`
	AltText    string    `gorm:"type:varchar(255)"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}
