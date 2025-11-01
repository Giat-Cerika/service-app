package models

import (
	"time"

	"github.com/google/uuid"
)

type Respondents struct {
	ID          uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID      uuid.UUID     `gorm:"type:uuid"`
	DetailID    uuid.UUID     `gorm:"type:uuid"`
	KuisionerID uuid.UUID     `gorm:"type:uuid"`
	Status      int           `gorm:"type:int" json:"status"`
	Duration    string        `gorm:"type:varchar(100)" json:"duration"`
	Answer      string        `gorm:"type:text" json:"answer"`
	Score       string        `gorm:"type:varchar(50)" json:"score"`
	User        User          `gorm:"foreignKey:UserID"`
	Detail      Detail        `gorm:"foreignKey:DetailID"`
	Kuisioner   Questionnaire `gorm:"foreignKey:KuisionerID"`
	CreatedAt   time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}
