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
	TransactionDate time.Time             `json:"transaction_date"`
	PaymentMethod    enums.PaymentMethod   `json:"payment_method"` // Add PaymentMethod
	PaymentCode      string                `json:"payment_code"`   // Add PaymentCode (if generated)
	PaymentStatus    enums.PaymentStatus   `json:"payment_status"` // Add PaymentStatus
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
}