package entities

import (
	"ApiRestFinance/internal/model/entities/enums"
	"gorm.io/gorm"
)

type LateFeeRule struct {
	gorm.Model
	EstablishmentID uint          `gorm:"index"`    // Optional, if rule is for a specific establishment
	Name            string        `gorm:"not null"` // E.g., "Standard Late Fee", "New Client Fee"
	DaysOverdueMin  int           `gorm:"not null"` // Minimum days overdue for rule to apply
	DaysOverdueMax  int           `gorm:"not null"` // Maximum days overdue
	FeeType         enums.FeeType `gorm:"not null"` // Percentage or fixed amount
	FeeValue        float64       `gorm:"not null"` // Percentage value or fixed amount
}
