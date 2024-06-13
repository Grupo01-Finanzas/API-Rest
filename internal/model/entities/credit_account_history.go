package entities

import (
	"ApiRestFinance/internal/model/entities/enums"
	"time"
)

type CreditAccountHistory struct {
	ID              uint                  `gorm:"primaryKey;autoIncrement"`
	CreditAccountID uint                  `gorm:"index;not null"`
	TransactionDate time.Time             `gorm:"not null"`
	TransactionType enums.TransactionType `gorm:"not null"`
	Amount          float64               `gorm:"not null"`  // Amount changed (positive or negative)
	Balance         float64               `gorm:"not null"`  // The resulting balance AFTER the event
	Description     string                `gorm:"type:text"` // Optional description
	CreditAccount   CreditAccount         `gorm:"foreignKey:CreditAccountID;references:ID"`
}
