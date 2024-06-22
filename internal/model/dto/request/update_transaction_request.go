package request

import "ApiRestFinance/internal/model/entities/enums"

type UpdateTransactionRequest struct {
	Amount          float64               `json:"amount" binding:"omitempty,gt=0"`
	Description     string                `json:"description" binding:"omitempty"`
	TransactionType enums.TransactionType `json:"transaction_type" binding:"omitempty"`
}