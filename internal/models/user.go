package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name        string     `gorm:"type:varchar(255)" json:"name"`
	Username    string     `gorm:"type:varchar(255);unique" json:"username"`
	Password    string     `gorm:"type:varchar(255)" json:"password"`
	Nisn        string     `gorm:"type:varchar(255)" json:"nisn"`
	DateOfBirth *time.Time `gorm:"type:date" json:"date_of_birth"`
	Age         int        `gorm:"type:int" json:"age"`
	Photo       *string    `gorm:"type:varchar(100)" json:"photo"`
	ClassID     uuid.UUID  `gorm:"type:uuid" json:"-"`
	Class       Class      `gorm:"foreignKey:ClassID" json:"class"`
	RoleID      uuid.UUID  `gorm:"type:uuid" json:"-"`
	Role        Role       `gorm:"foreignKey:RoleID" json:"role"`
	Status      int        `gorm:"type:int" json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}