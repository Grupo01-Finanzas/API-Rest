package request

import "ApiRestFinance/internal/model/entities/enums"

type UpdateLateFeeRuleRequest struct {
	Name           string        `json:"name" binding:"omitempty"`
	DaysOverdueMin int           `json:"days_overdue_min" binding:"omitempty,min=1"`
	DaysOverdueMax int           `json:"days_overdue_max" binding:"omitempty,gtfield=DaysOverdueMin"`
	FeeType        enums.FeeType `json:"fee_type" binding:"omitempty"`
	FeeValue       float64       `json:"fee_value" binding:"omitempty,gt=0"`
}
