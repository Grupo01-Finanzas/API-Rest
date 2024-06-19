package entities

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name            string          `gorm:"not null"`
	Description     string          `gorm:"not null"`
	Price           float64         `gorm:"not null"`
	Category        ProductCategory `gorm:"not null"`
	Stock           int             `gorm:"not null"`
	ImageUrl  		string     		`gorm:"default:'https://rahulindesign.websites.co.in/twenty-nineteen/img/defaults/product-default.png'"`
	IsActive        bool            `gorm:"not null"`
	CreatedAt       time.Time       `gorm:"not null"`
	UpdatedAt       time.Time       `gorm:"not null"`
	EstablishmentID uint            `gorm:"not null"`
	Establishment   Establishment   `gorm:"foreignKey:EstablishmentID;references:ID"`
}
