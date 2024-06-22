package entities

import (
	"ApiRestFinance/internal/model/entities/enums"
	"gorm.io/gorm"
	"time"
)

type Transaction struct {
	gorm.Model
	CreditAccountID  uint                   `gorm:"index;not null"`
	CreditAccount    *CreditAccount         `gorm:"foreignKey:CreditAccountID;references:ID"`
	TransactionType  enums.TransactionType `gorm:"not null"` // PURCHASE or PAYMENT
	Amount           float64               `gorm:"not null"`
	Description      string                `gorm:"type:text"`      // Optional description
	TransactionDate  time.Time             `gorm:"not null"`      // Date of the transaction
	PaymentMethod    enums.PaymentMethod   `gorm:"not null"`      // YAP, PLIN, CASH
	PaymentCode      string                `gorm:"default:null"`  // Code generated for client confirmation
	ConfirmationCode string                `gorm:"default:null"`  // Code provided by admin for confirmation
	PaymentStatus    enums.PaymentStatus   `gorm:"default:PENDING"` // PENDING, SUCCESS, FAILED
}