package request

import (
	"ApiRestFinance/internal/model/entities/enums"
)

type CreateTransactionRequest struct {
	CreditAccountID uint                  `json:"credit_account_id" binding:"required"`
	TransactionType enums.TransactionType `json:"transaction_type" binding:"required"`
	Amount          float64               `json:"amount" binding:"required,gt=0"`
	Description     string                `json:"description" binding:"omitempty"` // Optional
	RecipientType   enums.RecipientType   `json:"recipient_type" binding:"required"`
	RecipientID     uint                  `json:"recipient_id" binding:"required"`
}
