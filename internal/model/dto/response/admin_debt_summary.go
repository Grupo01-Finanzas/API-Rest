package response

import "time"

type AdminDebtSummary struct {
	ClientID       uint      `json:"client_id"`
	ClientName     string    `json:"client_name"`
	CreditType     string    `json:"credit_type"`
	InterestRate   float64   `json:"interest_rate"`
	NumberOfDues   int       `json:"number_of_installments"` // Only for long-term
	CurrentBalance float64   `json:"current_balance"`
	DueDate        time.Time `json:"due_date"` // For short-term or next installment
}
