package request

type CreateLateFeeRequest struct {
	CreditAccountID uint    `json:"credit_account_id" binding:"required"`
	Amount          float64 `json:"amount" binding:"required,gt=0"`
}
