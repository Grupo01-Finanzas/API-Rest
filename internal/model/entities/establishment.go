package entities

import (
	"gorm.io/gorm"
	"time"
)

type Establishment struct {
	gorm.Model
	RUC               string `gorm:"uniqueIndex;not null"`
	Name              string `gorm:"not null"`
	Phone             string `gorm:"not null"`
	Address           string `gorm:"not null"`
	ImageUrl          string `gorm:"default:'https://st2.depositphotos.com/47577860/46265/v/450/depositphotos_462652902-stock-illustration-building-business-company-icon.jpg'"`
	AdminID           uint
	Admin             *User     `gorm:"foreignKey:AdminID;references:ID"`
	IsActive          bool      `gorm:"not null"`
	LateFeePercentage float64   `gorm:"null"` // Added Late Fee Percentage
	CreatedAt         time.Time `gorm:"not null"`
	UpdatedAt         time.Time `gorm:"not null"`
}
