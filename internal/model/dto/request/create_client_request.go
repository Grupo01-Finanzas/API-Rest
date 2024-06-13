package request

type CreateClientRequest struct {
	UserID      uint    `json:"user_id" binding:"required"`
	Phone       string  `json:"phone" binding:"required"`
	Email       string  `json:"email" binding:"required,email"`
	CreditLimit float64 `json:"credit_limit" binding:"required,min=0"`
	IsActive    bool    `json:"is_active"`
}
