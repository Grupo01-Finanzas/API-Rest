package entities

import (
    "gorm.io/gorm"
    "time"
)

type LateFee struct {
    gorm.Model
    CreditAccountID uint       `gorm:"index;not null"`
    CreditAccount   CreditAccount `gorm:"foreignKey:CreditAccountID;references:ID"`
    Amount          float64    `gorm:"not null"`       // Amount of the late fee
    AppliedDate     time.Time  `gorm:"not null"`      // Date when the late fee was applied
}