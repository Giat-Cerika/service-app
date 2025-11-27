package models

import (
	"time"

	"github.com/google/uuid"
)

type ToootBrushLog struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID"`
	TimeType  string    `gorm:"type:varchar(100);not null" json:"time_type"`
	LogDate   time.Time `gorm:"type:date;not null" json:"log_date"`
	LogTime   time.Time `gorm:"type:time;not null" json:"log_time"`
	CreatedAt time.Time `json:"created_at"`
}
