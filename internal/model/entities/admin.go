package entities

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	UserID          uint           `gorm:"uniqueIndex;not null"`
	User            *User          `gorm:"foreignKey:UserID;references:ID"`
	EstablishmentID uint           `gorm:"uniqueIndex;not null"`
	Establishment   *Establishment `gorm:"foreignKey:EstablishmentID;references:ID"`
	IsActive        bool           `gorm:"not null"`
}
