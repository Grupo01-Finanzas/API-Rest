package request

import "ApiRestFinance/internal/model/entities/enums"

type CreateLateFeeRuleRequest struct {
	EstablishmentID uint          `json:"establishment_id" binding:"omitempty"` // Optional, if rule is for specific establishment
	Name            string        `json:"name" binding:"required"`
	DaysOverdueMin  int           `json:"days_overdue_min" binding:"required,min=1"`
	DaysOverdueMax  int           `json:"days_overdue_max" binding:"required,gtfield=DaysOverdueMin"`
	FeeType         enums.FeeType `json:"fee_type" binding:"required"`
	FeeValue        float64       `json:"fee_value" binding:"required,gt=0"`
}
