package response

import (
	"ApiRestFinance/internal/model/entities/enums"
	"time"
)

type CreditAccountResponse struct {
	ID                      uint                 `json:"id"`
	ClientID                uint                 `json:"client_id"`
	Client                  *UserResponse       `json:"client"`
	EstablishmentID         uint                 `json:"establishment_id"`
	Establishment           *EstablishmentResponse `json:"establishment"`
	CreditLimit             float64              `json:"credit_limit"`
	CurrentBalance          float64              `json:"current_balance"`
	MonthlyDueDate          int                  `json:"monthly_due_date"`
	InterestRate            float64              `json:"interest_rate"`
	InterestType            enums.InterestType   `json:"interest_type"`
	CreditType              enums.CreditType     `json:"credit_type"`
	GracePeriod             int                  `json:"grace_period"` 
	IsBlocked               bool                 `json:"is_blocked"`
	LastInterestAccrualDate time.Time            `json:"last_interest_accrual_date"`
	LateFeePercentage       float64            `json:"late_fee_percentage"`
	CreatedAt               time.Time            `json:"created_at"`
	UpdatedAt               time.Time            `json:"updated_at"`
}