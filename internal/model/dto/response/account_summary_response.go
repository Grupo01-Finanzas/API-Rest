package response

import "time"

// AccountSummaryResponse represents a summary of a client's account.
type AccountSummaryResponse struct {
	CurrentBalance float64               `json:"current_balance"`
	DueDate        time.Time             `json:"due_date"`
	TotalInterest  float64               `json:"total_interest"`
	Transactions   []TransactionResponse `json:"transactions"`
}
