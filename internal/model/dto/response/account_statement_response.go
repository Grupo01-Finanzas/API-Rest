package response

import "time"

// AccountStatementResponse defines the response structure for a client account statement.
type AccountStatementResponse struct {
    ClientID        uint                  `json:"client_id"`
    StartDate       time.Time             `json:"start_date"`
    EndDate         time.Time             `json:"end_date"`
    StartingBalance float64               `json:"starting_balance"`
    Transactions    []TransactionResponse `json:"transactions"`
}