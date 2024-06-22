package response

import "time"

type LateFeeResponse struct {
	ID              uint      `json:"id"`
	CreditAccountID uint      `json:"credit_account_id"`
	Amount          float64   `json:"amount"`
	AppliedDate     time.Time `json:"applied_date"`
}