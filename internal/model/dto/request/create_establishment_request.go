package request

type CreateEstablishmentRequest struct {
	RUC               string  `json:"ruc" binding:"required"`
	Name              string  `json:"name" binding:"required"`
	Phone             string  `json:"phone" binding:"required"`
	Address           string  `json:"address" binding:"required"`
	ImageUrl          string  `json:"image_url" binding:"omitempty"`
	LateFeePercentage float64 `json:"late_fee_percentage" binding:"omitempty"`
}
