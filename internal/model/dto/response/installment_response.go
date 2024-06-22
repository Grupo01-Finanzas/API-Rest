package response

import (
	"ApiRestFinance/internal/model/entities/enums"
	"time"
)

type InstallmentResponse struct {
	ID              uint                    `json:"id"`
	CreditAccountID uint                    `json:"credit_account_id"`
	DueDate         time.Time               `json:"due_date"`
	Amount          float64                 `json:"amount"`
	Status          enums.InstallmentStatus `json:"status"`
	CreatedAt       time.Time               `json:"created_at"`
	UpdatedAt       time.Time               `json:"updated_at"`
}