package entities

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	UserID      uint    `gorm:"uniqueIndex;not null"`
	User        User    `gorm:"foreignKey:UserID;references:ID"`
	Phone       string  `gorm:"not null"`
	Email       string  `gorm:"uniqueIndex;not null"`
	CreditLimit float64 `gorm:"not null"`
	IsActive    bool    `gorm:"not null"`
}
