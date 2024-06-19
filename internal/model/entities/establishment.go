package entities

import "gorm.io/gorm"

type Establishment struct {
	gorm.Model
	RUC           string        `gorm:"uniqueIndex;not null"`
	Name          string        `gorm:"uniqueIndex;not null"`
	Phone         string        `gorm:"not null"`
	Address       string        `gorm:"not null"`
	ImageUrl  	  string     	`gorm:"default:'https://st2.depositphotos.com/47577860/46265/v/450/depositphotos_462652902-stock-illustration-building-business-company-icon.jpg'"`
	Admin         *Admin        `gorm:"foreignKey:EstablishmentID;references:ID"`
	IsActive      bool          `gorm:"not null"`
	Clients       []Client      `gorm:"many2many:establishment_clients;"`
	Products      []Product     `gorm:"foreignKey:EstablishmentID;references:ID"`
	LateFeeRuleID uint          `gorm:"uniqueIndex;not null"`
	LateFeeRules  []LateFeeRule `gorm:"foreignKey:EstablishmentID;references:ID"` // Relationship for late fee rules
}
