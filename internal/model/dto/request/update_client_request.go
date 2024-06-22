package request

import (
	"ApiRestFinance/internal/model/entities/enums"
)

// UpdateClientRequest represents the request to update an existing client.
type UpdateClientRequest struct {
	Name             string             `json:"name" binding:"omitempty"`
	Address          string             `json:"address" binding:"omitempty"`
	Phone            string             `json:"phone" binding:"omitempty"`
	IsActive         bool               `json:"is_active"`
	CreditLimit     float64            `json:"credit_limit" binding:"omitempty,gt=0"`
	MonthlyDueDate  int                `json:"monthly_due_date" binding:"omitempty,min=1,max=31"`
	InterestRate    float64            `json:"interest_rate" binding:"omitempty,gt=0"`
	InterestType    enums.InterestType `json:"interest_type" binding:"omitempty"`
	CreditType      enums.CreditType   `json:"credit_type" binding:"omitempty"`
	GracePeriod     int                `json:"grace_period" binding:"omitempty,min=0"`
}