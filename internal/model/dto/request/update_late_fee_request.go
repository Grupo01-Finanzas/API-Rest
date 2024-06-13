package request

type UpdateLateFeeRequest struct {
	Amount float64 `json:"amount" binding:"omitempty,gt=0"`
}
