package request

import (
	"ApiRestFinance/internal/model/entities/enums"
	"time"
)

type UpdateInstallmentRequest struct {
	DueDate time.Time               `json:"due_date" binding:"omitempty"`
	Amount  float64                 `json:"amount" binding:"omitempty,gt=0"`
	Status  enums.InstallmentStatus `json:"status" binding:"omitempty"`
}
