package response

// ClientBalanceResponse represents the response for the client's balance
type ClientBalanceResponse struct {
	ClientID       uint    `json:"client_id"`
	CurrentBalance float64 `json:"current_balance"`
}
