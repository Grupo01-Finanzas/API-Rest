package entities

import "gorm.io/gorm"

type Establishment struct {
	gorm.Model
	RUC           string        `gorm:"uniqueIndex;not null"`
	Name          string        `gorm:"uniqueIndex;not null"`
	Phone         string        `gorm:"not null"`
	Address       string        `gorm:"not null"`
	Admin         *Admin        `gorm:"foreignKey:EstablishmentID;references:ID"`
	IsActive      bool          `gorm:"not null"`
	Clients       []Client      `gorm:"many2many:establishment_clients;"`
	Products      []Product     `gorm:"foreignKey:EstablishmentID;references:ID"`
	LateFeeRuleID uint          `gorm:"uniqueIndex;not null"`
	LateFeeRules  []LateFeeRule `gorm:"foreignKey:EstablishmentID;references:ID"` // Relationship for late fee rules
}
