package request

import (
	"ApiRestFinance/internal/model/entities/enums"
)

type CreateCreditAccountRequest struct {
	ClientID        uint               `json:"client_id" binding:"required"`
	CreditLimit     float64            `json:"credit_limit" binding:"required,gt=0"`
	MonthlyDueDate  int                `json:"monthly_due_date" binding:"required,min=1,max=31"`
	InterestRate    float64            `json:"interest_rate" binding:"required,gt=0"`
	InterestType    enums.InterestType `json:"interest_type" binding:"required"`
	CreditType      enums.CreditType   `json:"credit_type" binding:"required"`
	GracePeriod     int                `json:"grace_period" binding:"omitempty,min=0"` // Optional, for long-term credit
}