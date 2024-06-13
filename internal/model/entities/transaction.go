package entities

import (
	"ApiRestFinance/internal/model/entities/enums"
	"gorm.io/gorm"
	"time"
)

type Transaction struct {
	gorm.Model
	CreditAccountID uint                  `gorm:"index"`    // Now optional
	RecipientType   enums.RecipientType   `gorm:"not null"` // New field to indicate recipient type
	RecipientID     uint                  `gorm:"not null"` // ID of the recipient
	TransactionType enums.TransactionType `gorm:"not null"`
	Amount          float64               `gorm:"not null"`
	Description     string                `gorm:"type:text"` // Optional description
	TransactionDate time.Time             `gorm:"not null"`  // Date of the transaction
	CreditAccount   CreditAccount         `gorm:"foreignKey:CreditAccountID;references:ID"`
}
