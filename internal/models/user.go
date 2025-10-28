package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name        *string    `gorm:"type:varchar(255)" json:"name"`
	Username    string     `gorm:"type:varchar(255);index" json:"username"`
	Password    string     `gorm:"type:varchar(255)" json:"password"`
	Nisn        *string    `gorm:"type:varchar(255);index" json:"nisn"`
	DateOfBirth *time.Time `gorm:"type:date;" json:"date_of_birth"`
	Age         *time.Time `gorm:"type:int;" json:"age"`
	Photo       string     `gorm:"type:varchar(255);index" json:"photo"`
	RoleID      uuid.UUID  `gorm:"type:uuid"`
	ClassID     *uuid.UUID `gorm:"type:uuid"`
	Role        Role       `gorm:"foreignKey:RoleID"`
	Class       Class      `gorm:"foreignKey:ClassID"`
	Status      int        `gorm:"type:int;" json:"status"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}
