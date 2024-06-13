package response

import (
	"ApiRestFinance/internal/model/entities/enums"
	"time"
)

type TransactionResponse struct {
	ID              uint                  `json:"id"`
	CreditAccountID uint                  `json:"credit_account_id"`
	TransactionType enums.TransactionType `json:"transaction_type"`
	Amount          float64               `json:"amount"`
	Description     string                `json:"description"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
}
