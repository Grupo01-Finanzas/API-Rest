package request

import "time"

type UpdateCreditRequestDueDate struct {
    DueDate time.Time `json:"due_date" binding:"required"`
}