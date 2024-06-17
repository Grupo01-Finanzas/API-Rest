package entities

import (
	"gorm.io/gorm"
	"time"
)

type Client struct {
	gorm.Model
	UserID    uint      `gorm:"uniqueIndex;not null"`
	User      *User     `gorm:"foreignKey:UserID;references:ID"`
	IsActive  bool      `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}
