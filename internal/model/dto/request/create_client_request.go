package request

import (
	"ApiRestFinance/internal/model/entities/enums"
)

// CreateClientRequest represents the request to create a new client.
type CreateClientRequest struct {
	EstablishmentID   uint               `json:"establishment_id" binding:"required"`
	DNI               string             `json:"dni" binding:"required,min=8,max=8"`
	Email             string             `json:"email" binding:"omitempty,email"` // Optional email
	Name              string             `json:"name" binding:"required"`
	Address           string             `json:"address" binding:"required,min=5"`
	Phone             string             `json:"phone" binding:"required,min=9,max=9"`
	CreditLimit       float64            `json:"credit_limit" binding:"required,gt=0"`
	MonthlyDueDate    int                `json:"monthly_due_date" binding:"required,min=1,max=31"`
	InterestRate      float64            `json:"interest_rate" binding:"required,gt=0.0"`
	InterestType      enums.InterestType `json:"interest_type" binding:"required"`
	CreditType        enums.CreditType   `json:"credit_type" binding:"required"`
	GracePeriod       int                `json:"grace_period" binding:"omitempty,min=0"`
	LateFeePercentage float64            `json:"late_fee_percentage" binding:"omitempty"`
}
