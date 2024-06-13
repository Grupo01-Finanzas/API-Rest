package entities

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Name     string
	Roles    []Role  `gorm:"many2many:user_roles;"`
	Client   *Client `gorm:"foreignKey:UserID;references:ID"`
	Admin    *Admin  `gorm:"foreignKey:UserID;references:ID"`
}
