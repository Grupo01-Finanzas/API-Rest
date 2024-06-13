package entities

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name  string `gorm:"uniqueIndex;not null"`
	Users []User `gorm:"many2many:user_roles;"`
}
