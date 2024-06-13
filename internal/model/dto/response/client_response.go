package response

import "time"

type ClientResponse struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	CreditLimit float64   `json:"credit_limit"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
