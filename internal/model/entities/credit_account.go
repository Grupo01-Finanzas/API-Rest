package entities

import (
    "ApiRestFinance/internal/model/entities/enums"
    "gorm.io/gorm"
    "time"
)

type CreditAccount struct {
	gorm.Model
	ClientID                uint               `gorm:"index;not null"`
	Client                  *User            `gorm:"foreignKey:ClientID;references:ID"` // Client this account belongs to
	EstablishmentID         uint               `gorm:"index;not null"`
	Establishment           *Establishment     `gorm:"foreignKey:EstablishmentID;references:ID"`
	CreditLimit             float64            `gorm:"not null"`
	CurrentBalance          float64            `gorm:"not null"` // Current balance owed
	MonthlyDueDate          int                `gorm:"not null"` // Day of the month (1-31) when payment is due
	InterestRate            float64            `gorm:"not null"` // Annual interest rate
	InterestType            enums.InterestType `gorm:"not null"` // NOMINAL or EFFECTIVE
	CreditType              enums.CreditType   `gorm:"not null"` // SHORT_TERM or LONG_TERM
	GracePeriod             int                `gorm:"default:0"` // Grace period in months (for LONG_TERM credit)
	IsBlocked               bool               `gorm:"default:false"`
	LastInterestAccrualDate time.Time          `gorm:"not null"` // Date when interest was last applied
	LateFeePercentage       float64            `gorm:"not null"` // Percentage for late fee calculation
	CreatedAt               time.Time          `gorm:"not null"`
	UpdatedAt               time.Time          `gorm:"not null"`
}