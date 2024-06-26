package entities

import (
	"ApiRestFinance/internal/model/entities/enums"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	DNI       string     `gorm:"uniqueIndex;not null"`
	Email     string     `gorm:"uniqueIndex;not null"` 
	Password  string     `gorm:"not null"`
	Name      string     `gorm:"not null"`
	Address   string     `gorm:"not null"`
	Phone     string     `gorm:"not null"`
	PhotoUrl  string     `gorm:"default:'https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png'"`
	Rol       enums.Role `gorm:"type:text;not null"` // ADMIN or CLIENT
	CreatedAt time.Time  `gorm:"not null"`
	UpdatedAt time.Time  `gorm:"not null"`
}