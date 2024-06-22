package entities

import (
	"ApiRestFinance/internal/model/entities/enums"
	"gorm.io/gorm"
	"time"
)

type Installment struct {
	gorm.Model
	CreditAccountID uint                    `gorm:"index;not null"`
	CreditAccount   *CreditAccount           `gorm:"foreignKey:CreditAccountID;references:ID"`
	DueDate         time.Time               `gorm:"not null"` // Due date of the installment
	Amount          float64                 `gorm:"not null"`
	Status          enums.InstallmentStatus `gorm:"not null;default:PENDING"` // PENDING, PAID, OVERDUE
}