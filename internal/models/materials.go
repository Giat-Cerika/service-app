package models

import (
	"time"

	"github.com/google/uuid"
)

type Materials struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Title       string    `gorm:"type:varchar(255);index" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedBy   uuid.UUID `gorm:"type:uuid"`
	User        User      `gorm:"foreignKey:CreatedBy"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
