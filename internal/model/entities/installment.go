package entities

import (
	"ApiRestFinance/internal/model/entities/enums"
	"gorm.io/gorm"
	"time"
)

type Installment struct {
	gorm.Model
	CreditAccountID uint                    `gorm:"index;not null"`
	DueDate         time.Time               `gorm:"not null"`
	Amount          float64                 `gorm:"not null"`
	Status          enums.InstallmentStatus `gorm:"not null"`
	CreditAccount   CreditAccount           `gorm:"foreignKey:CreditAccountID;references:ID"`
}
