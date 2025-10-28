package models

import (
	"time"

	"github.com/google/uuid"
)

type Class struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	NameClass string    `gorm:"type:varchar(255);index" json:"name_class"`
	Grade     string    `gorm:"type:varchar(255);index" json:"grade"`
	Teacher   string    `gorm:"type:varchar(255);index" json:"teacher"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
