package response

import (
	"ApiRestFinance/internal/model/entities/enums"
	"time"
)

type CreditAccountResponse struct {
	ID                      uint                 `json:"id"`
	EstablishmentID         uint                 `json:"establishment_id"`
	ClientID                uint                 `json:"client_id"`
	CreditLimit             float64              `json:"credit_limit"`
	CurrentBalance          float64              `json:"current_balance"`
	MonthlyDueDate          int                  `json:"monthly_due_date"`
	InterestRate            float64              `json:"interest_rate"`
	InterestType            enums.InterestType   `json:"interest_type"`
	CreditType              enums.CreditType     `json:"credit_type"`
	GracePeriod             int                  `json:"grace_period"`
	IsBlocked               bool                 `json:"is_blocked"`
	LastInterestAccrualDate time.Time            `json:"last_interest_accrual_date"`
	CreatedAt               time.Time            `json:"created_at"`
	UpdatedAt               time.Time            `json:"updated_at"`
	LateFeeRuleID           uint                 `json:"late_fee_rule_id"`
	Client                  *ClientResponse      `json:"client"`
	LateFeeRule             *LateFeeRuleResponse `json:"late_fee_rule"`
}
