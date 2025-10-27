package models

import (
	"time"

	"github.com/google/uuid"
)

type Material struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Title       string    `gorm:"type:varchar(100);not null" json:"title"`
	Description string    `gorm:"type:text;not null" json:"description"`
	CreatedBy   uuid.UUID `gorm:"type:uuid; not null" json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
