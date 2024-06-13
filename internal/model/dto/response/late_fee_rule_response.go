package response

import "ApiRestFinance/internal/model/entities/enums"

type LateFeeRuleResponse struct {
	ID              uint          `json:"id"`
	EstablishmentID uint          `json:"establishment_id"`
	Name            string        `json:"name"`
	DaysOverdueMin  int           `json:"days_overdue_min"`
	DaysOverdueMax  int           `json:"days_overdue_max"`
	FeeType         enums.FeeType `json:"fee_type"`
	FeeValue        float64       `json:"fee_value"`
}
