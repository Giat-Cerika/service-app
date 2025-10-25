package models

type Role struct {
    BaseModel
    Name     string `gorm:"type:varchar(100);not null" json:"name"`
    Email    string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
    Password string `gorm:"type:varchar(255);not null" json:"-"`
    Phone    string `gorm:"type:varchar(20)" json:"phone"`
    NIM      string `gorm:"type:varchar(50);uniqueIndex" json:"nim"`
    Major    string `gorm:"type:varchar(100)" json:"major"`
}