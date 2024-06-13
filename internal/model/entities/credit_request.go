package entities

import (
	"ApiRestFinance/internal/model/entities/enums"
	"gorm.io/gorm"
	"time"
)

type CreditRequest struct {
	gorm.Model
	ClientID             uint                `gorm:"index;not null"`
	EstablishmentID      uint                `gorm:"index;not null"`
	RequestedCreditLimit float64             `gorm:"not null"`
	MonthlyDueDate       int                 `gorm:"not null"` // Day of the month (1-31)
	InterestType         enums.InterestType  `gorm:"not null"`
	InterestRate         float64             `gorm:"not null"`
	CreditType           enums.CreditType    `gorm:"not null"`
	GracePeriod          int                 `gorm:"default:0"` // In months
	Status               CreditRequestStatus `gorm:"not null;default:PENDING"`
	ApprovedAt           *time.Time          `gorm:"null"` // Timestamp when the request was approved
	RejectedAt           *time.Time          `gorm:"null"` // Timestamp when the request was rejected
	Client               Client              `gorm:"foreignKey:ClientID;references:ID"`
	Establishment        Establishment       `gorm:"foreignKey:EstablishmentID;references:ID"`
}
