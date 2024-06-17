package entities

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	UserID   uint  `gorm:"uniqueIndex;not null"`
	User     *User `gorm:"foreignKey:UserID;references:ID"`
	IsActive bool  `gorm:"not null"`
}
