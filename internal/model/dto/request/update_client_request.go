package request

type UpdateClientRequest struct {
	Phone       string  `json:"phone" binding:"required"`
	Email       string  `json:"email" binding:"required,email"`
	CreditLimit float64 `json:"credit_limit" binding:"required,min=0"`
	IsActive    bool    `json:"is_active"`
}
