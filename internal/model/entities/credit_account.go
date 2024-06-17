package entities

import (
	"ApiRestFinance/internal/model/entities/enums"
	"gorm.io/gorm"
	"time"
)

// CreditAccount represents a client's credit account.
type CreditAccount struct {
	gorm.Model
	EstablishmentID         uint               `gorm:"index;not null"`
	ClientID                uint               `gorm:"index;not null"`
	CreditLimit             float64            `gorm:"not null"`
	MonthlyDueDate          int                `gorm:"not null"` // Day of the month (1-31)
	InterestRate            float64            `gorm:"not null"`
	InterestType            enums.InterestType `gorm:"not null"`
	CreditType              enums.CreditType   `gorm:"not null"`
	GracePeriod             int                `gorm:"default:0"` // In months
	IsBlocked               bool               `gorm:"default:false"`
	LastInterestAccrualDate time.Time          `gorm:"not null"`
	CurrentBalance          float64            `gorm:"not null"`
	Establishment           *Establishment     `gorm:"foreignKey:EstablishmentID;references:ID"`
	Client                  *Client            `gorm:"foreignKey:ClientID;references:ID"`
	Transactions            []*Transaction     `gorm:"foreignKey:CreditAccountID;references:ID"`
	LateFees                []*LateFee         `gorm:"foreignKey:CreditAccountID;references:ID"` // New relationship
	Installments            []*Installment     `gorm:"foreignKey:CreditAccountID;references:ID"` // New relationship (for long-term credit)
	LateFeeRuleID           uint               `gorm:"index"`                                    // Foreign key to LateFeeRule
	LateFeeRule             *LateFeeRule       `gorm:"foreignKey:LateFeeRuleID;references:ID"`   // Relationship with LateFeeRule
}
