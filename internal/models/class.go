package models

import (
	"time"

	"github.com/google/uuid"
)

type Class struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	NameClass string    `gorm:"type:varchar(100);not null" json:"name_class"`
	Level     int       `gorm:"not null" json:"level"`
	HTeacher  string    `grom:"type:varchar(200);not null" json:"h_teacher"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
