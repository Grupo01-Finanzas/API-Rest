package entities

import (
	"gorm.io/gorm"
	"time"
)

type LateFee struct {
	gorm.Model
	CreditAccountID uint          `gorm:"index;not null"`
	Amount          float64       `gorm:"not null"`
	AppliedDate     time.Time     `gorm:"not null"`
	CreditAccount   CreditAccount `gorm:"foreignKey:CreditAccountID;references:ID"`
}
