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
	Rol       enums.Role `gorm:"type:text;not null"`
	Client    *Client    `gorm:"foreignKey:UserID;references:ID"`
	Admin     *Admin     `gorm:"foreignKey:UserID;references:ID"`
	CreatedAt time.Time  `gorm:"not null"`
	UpdatedAt time.Time  `gorm:"not null"`
}
