package response

import (
	"ApiRestFinance/internal/model/entities"
	"ApiRestFinance/internal/model/entities/enums"
	"time"
)

type CreditRequestResponse struct {
	ID                   uint                         `json:"id"`
	ClientID             uint                         `json:"client_id"`
	EstablishmentID      uint                         `json:"establishment_id"`
	RequestedCreditLimit float64                      `json:"requested_credit_limit"`
	MonthlyDueDate       int                          `json:"monthly_due_date"`
	InterestType         enums.InterestType           `json:"interest_type"`
	CreditType           enums.CreditType             `json:"credit_type"`
	GracePeriod          int                          `json:"grace_period"`
	Status               entities.CreditRequestStatus `json:"status"`
	ApprovedAt           *time.Time                   `json:"approved_at"`
	RejectedAt           *time.Time                   `json:"rejected_at"`
	CreatedAt            time.Time                    `json:"created_at"`
	UpdatedAt            time.Time                    `json:"updated_at"`
}
